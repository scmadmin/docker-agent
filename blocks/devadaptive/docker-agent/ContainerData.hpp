//
// Created by paul on 6/7/15.
//

#ifndef DOCKER_AGENT_CONTAINERDATA_H
#define DOCKER_AGENT_CONTAINERDATA_H

#include "stdlib.h"
#include <map>
#include <string>
#include "hjiang/jsonxx/jsonxx.h"
#include "dautil.hpp"

class ContainerData {

public:

    ContainerData(const std::string& id, const std::string& tenantId) : containerId(id), tenantId(tenantId) {}

    void putMetricData(const std::string& key, const long value) {
        metricData.insert({key, value});
    }

    void putDocumentData(const std::string& key, const std::string& value) {
        documentData.insert({key, value});
    }

    //this should be external
    std::string metricToJSON() {
        jsonxx::Object metricObject;
        metricObject << "containerId" << containerId;
        metricObject << "agentId" << "xyz-agentid";
        for (auto metricPair : metricData) {
            metricObject << metricPair.first << metricPair.second;
        }
        return metricObject.json();
    }

    //this should be external

    jsonxx::Array getMetricArray() {
        jsonxx::Array metricArray;
        long now = currentTimeMillis();
        for (auto metricPair : metricData) {
            jsonxx::Object metricObject;
            metricObject << "metricId1" << metricPair.first;
            metricObject << "metricId2" << containerId;
            metricObject << "agentId" << "xyz-agentid";
            metricObject << "tenantId" << tenantId;
            metricObject << "day" << 0;
            metricObject << "duration" << duration;
            metricObject << "eventTime" << now;
            metricObject << "metricType" << "not used";
            jsonxx::Object dataObject;
            dataObject << "value" << metricPair.second;
            dataObject << "count" << 1;
            metricObject << "data" << dataObject;
            metricArray  << metricObject;
        }
        return metricArray;
    }

    std::string documentToJSON() {
        jsonxx::Object documentObject;
        documentObject << "containerId" << containerId;
        documentObject << "agentId"  << "xyz-agentid";
        for(auto documentPair : documentData)  {
            documentObject << documentPair.first <<  documentPair.second;
        }
        return documentObject.json();
    }

    void dump() {
        std::cout << "container: " << containerId << endl;
        std::cout << "metric data" << endl;
        for(auto metricPair : metricData)  {
            std::cout << metricPair.first << " " << metricPair.second << endl;
        }

        std::cout << "document data" << endl;
        for(auto documentPair : documentData)  {
            std::cout << documentPair.first << " " << documentPair.second << endl;
        }
    }


private:
    std::map<const std::string, const std::string> documentData;
    std::map<const std::string, const long> metricData;
    const std::string containerId;
    const std::string tenantId;
    long duration = 15000L;

};


#endif //DOCKER_AGENT_CONTAINERDATA_H


