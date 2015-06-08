//
// Created by paul on 6/7/15.
//

#ifndef DOCKER_AGENT_CONTAINERDATA_H
#define DOCKER_AGENT_CONTAINERDATA_H

#include "stdlib.h"
#include <map>
#include <string>
#include "hjiang/jsonxx/jsonxx.h"

class ContainerData {

public:

    ContainerData(std::string& id) : containerId(id) {}

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
    std::string containerId;


};


#endif //DOCKER_AGENT_CONTAINERDATA_H
