//
// Created by paul on 6/7/15.
//

#ifndef DOCKER_AGENT_CONTAINERDATA_H
#define DOCKER_AGENT_CONTAINERDATA_H

#include "stdlib.h"
#include <map>
#include <string>
#include "dautil.hpp"

using MetricGroupAndKey = std::pair<const std::string, const std::string>;


class ContainerData {

public:

    ContainerData(const std::string& id, const std::string& tenantId, const std::string& hostname) :
            containerId(id), tenantId(tenantId), hostname(hostname) {}

    void putMetricData(const std::string& metricGroup, const std::string& metricKey, const long value) {
        MetricGroupAndKey groupAndKey(metricGroup, metricKey);
        metricData.insert({groupAndKey, value});
    }


    void putDocumentData(const std::string& key, const std::string& value) {
        documentData.insert({key, value});
    }

    template <typename Func>
    void mapMetricArray(Func mapper) {
        long now = currentTimeMillis();
        for (auto metricPair : metricData) {
            mapper(metricPair, this, now);
        }
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


    const string &getContainerId() const {
        return containerId;
    }

    const string &getTenantId() const {
        return tenantId;
    }

    const string &getHostname() const {
        return hostname;
    }

    long getDuration() const {
        return duration;
    }


    const string &getName() const {
        return name;
    }

    void setName(const string &name) {
        ContainerData::name = name;
    }

private:
    std::map<const std::string, const std::string> documentData;
    std::map<MetricGroupAndKey, const long> metricData;
    const std::string containerId;
    const std::string tenantId;
    const std::string hostname;
    std::string name;
    long duration = 15000L;

};


#endif //DOCKER_AGENT_CONTAINERDATA_H


