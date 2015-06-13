//
// Created by paul on 6/10/15.
//

#ifndef DOCKER_AGENT_HTTPCLIENT_H
#define DOCKER_AGENT_HTTPCLIENT_H

#include "hjiang/jsonxx/jsonxx.h"

//#include "lasote/lambda_http_client/http_request.h"

#include <string.h>
#include <iostream>
#include <functional>
#include <fstream>
#include <memory>

using namespace httpmodels;
using namespace lasote;


class HTTPClient {
    std::string uri;
    //https://devadaptive.com:/p/porter



    shared_ptr<LambdaRequest> post_and_parse_json_response(string path, string body, JSON_F ret_cb, ERR_F err_cb, CHAR_PTR_F on_body_cb=CHAR_PTR_F()){
        log_info(path);
        httpmodels::Request request;
        httpmodels::Method method("POST", uri + path);
        request.method = &method;
        request.body = body;
        std::pair<string,string> content_type("Content-Type", "application/json");
        request.headers.insert(content_type);

        ostringstream convert;
        convert << body.length();

        std::pair<string,string> content_len("Content-Length", convert.str());
        request.headers.insert(content_len);

        return call_and_parse_response(request, ret_cb, err_cb, on_body_cb);
    }

    shared_ptr<LambdaRequest> call_and_parse_response(Request& request, JSON_F ret_cb, ERR_F err_cb, CHAR_PTR_F on_body_cb){
        auto request_call = std::make_shared<LambdaRequest>();

        request_call->on_message_complete_cb = [request_call, ret_cb, err_cb] (int status) {
            if(status > 299){
                log_debug("Error calling: " << status);
                if(err_cb != NULL){
                    err_cb(status, request_call->response_buffer);
                }
            }
            else{
                log_debug("Status: " << status << endl);
                JSON_OBJECT json_ret;
                if(request_call->response_buffer.length() > 0){
                    if(request_call->response_buffer[0] == '{'){
                        json_ret.parse(request_call->response_buffer);
                        log_debug("Response JSON OBJECT: " << json_ret << endl);
                    }
                    else if(request_call->response_buffer[0] == '['){
                        JSON_ARRAY json_tmp;
                        json_tmp.parse(request_call->response_buffer);
                        json_ret << "data" << json_tmp;
                        log_debug("Response JSON ARRAY: " << json_tmp << endl);
                    }
                }
                ret_cb(json_ret);
            }
        };
        request_call->on_body_cb = on_body_cb;
        request_call->send(request);
        return request_call;
    }

public:
    HTTPClient(string host) : uri(host){
    }

    shared_ptr<LambdaRequest> postMetrics(jsonxx::Array metrics, JSON_F ret_cb, ERR_F err_cb) {
        std::string path("/metrics");
        return post_and_parse_json_response(path, metrics.json(), ret_cb, err_cb);
    }
};


#endif //DOCKER_AGENT_HTTPCLIENT_H
