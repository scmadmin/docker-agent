//
// Created by root on 6/13/15.
//

#ifndef DOCKER_AGENT_CURL_CLIENT_HPP
#define DOCKER_AGENT_CURL_CLIENT_HPP

#include <string>
#include <curl/curl.h>

bool isCurlInittedx = false;

static size_t writeCallback(void *contents, size_t size, size_t nmemb, void *userp) {
    //std::cout << "debug: " << contents << std::endl;
    ((std::string *) userp)->append((char *) contents, size * nmemb);
    return size * nmemb;
}


int postJSON(std::string& url, std::string& json, std::string& responseBuffer) {
    CURL *curl;
    CURLcode res;

    if (!isCurlInittedx) {
        curl_global_init(CURL_GLOBAL_ALL);
        isCurlInittedx = true;
    }

    curl = curl_easy_init();
    if (curl) {
        struct curl_slist *headers = NULL;
        headers = curl_slist_append(headers, "Accept: application/json");
        headers = curl_slist_append(headers, "Content-Type: application/json");
        headers = curl_slist_append(headers, "charsets: utf-8");
        //curl_easy_setopt(curl, CURLOPT_VERBOSE, 1L);
        //curl_easy_setopt(curl, CURLOPT_HEADER, 1L);
        curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, json.c_str());
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, writeCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &responseBuffer);

        //todo:  to keep curl from posting the result to stdout, assign a function to read the data
        // from the response.

        //curl_easy_setopt(curl, CURLOPT_POSTFIELDSIZE, json.length());

        res = curl_easy_perform(curl);
        if (res != CURLE_OK) {
            std::cout << "curl_easy_perform() failed: " << curl_easy_strerror(res) << std::endl;
            return -1;
        }
        curl_easy_cleanup(curl);

        //std::cout << "response: " << readBuffer << std::endl;
    } else {
        throw "couldn't obtain curl handle";
    }
    //curl_global_cleanup();
    return 0;
}

#endif //DOCKER_AGENT_CURL_CLIENT_HPP
