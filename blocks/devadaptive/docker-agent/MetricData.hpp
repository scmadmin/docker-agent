//
// Created by paul on 6/30/15.
//

#ifndef DOCKER_AGENT_METRICDATA_HPP
#define DOCKER_AGENT_METRICDATA_HPP

#include <string>

class MetricData {


private:
    std::string metricId;
    std::string metricGroup;
    long value;

};
#endif //DOCKER_AGENT_METRICDATA_HPP
