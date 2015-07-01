//
// Created by paul on 6/14/15.
//

#ifndef DOCKER_AGENT_JSONMETRIC_HPP
#define DOCKER_AGENT_JSONMETRIC_HPP

#include "../../../deps/hjiang/jsonxx/jsonxx.h"

using MetricMapperType  = std::function<void(std::pair<const std::string &, const long> metricPair,
                                             const std::string &containerId, long now, long duration,
                                             const std::string &tenantId)>;

class JSONMetrics {

public:
    jsonxx::Array &getMetricsArray() {
        return jsonMetrics;
    }

    MetricMapperType &getMapper() {
        return mapper;
    }

    int getContainerCount() {
        return containerCount;
    }

    void incrementContainerCount() {
        containerCount++;
    }

private:
    jsonxx::Array jsonMetrics;
    int containerCount = 0;

    MetricMapperType mapper = [&](std::pair<const std::string &, const long> metricPair,
                                  const std::string &containerId, long now, long duration,
                                  const std::string &tenantId) {
        jsonxx::Object metricObject;
        metricObject << "host" << "192.168.0.1";
        metricObject << "container" << containerId;
        metricObject << "containerName" << "TODO";
        metricObject << "metricId" << metricPair.first;
        metricObject << "metricType" << "TODO e.g. memory";
        metricObject << "agentId" << "c4e99f12-666b-a666-97b8-6e666db9a667";
        metricObject << "tenantId" << tenantId;
        metricObject << "duration" << duration;
        metricObject << "eventTime" << now;
        metricObject << "dataSource" << "docker";
        jsonxx::Object dataObject;
        dataObject << "value" << metricPair.second;
        dataObject << "count" << 1;
        metricObject << "data" << dataObject;
        jsonMetrics << metricObject;
    };
};

#endif //DOCKER_AGENT_JSONMETRIC_HPP


