//
// Created by paul on 6/14/15.
//

#ifndef DOCKER_AGENT_JSONMETRIC_HPP
#define DOCKER_AGENT_JSONMETRIC_HPP

#include "../../../deps/hjiang/jsonxx/jsonxx.h"
#include "ContainerData.hpp"

using MetricMapperType  = std::function<void(std::pair<const MetricGroupAndKey &, const long> metricPair,
                                             const ContainerData* containerData, const long now)>;

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

    MetricMapperType mapper = [&](std::pair<const MetricGroupAndKey &, const long> metricPair,
                                  const ContainerData* containerData, const long now) {
        jsonxx::Object metricObject;
        metricObject << "host" <<  containerData->getHostname();
        metricObject << "container" << containerData->getContainerId();
        metricObject << "containerName" << containerData->getName();
        metricObject << "metricId" << metricPair.first.second;
        metricObject << "metricGroup" << metricPair.first.first;
        metricObject << "agentId" << "c4e99f12-666b-a666-97b8-6e666db9a667";
        metricObject << "tenantId" << containerData->getTenantId();
        metricObject << "duration" << containerData->getDuration();
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


