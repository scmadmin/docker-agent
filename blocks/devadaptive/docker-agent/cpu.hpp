//
// Created by paul on 6/7/15.
//

#ifndef DOCKER_AGENT_CPU_HPP
#define DOCKER_AGENT_CPU_HPP



#include <string>
#include <map>
#include <fstream>


void getCPUData(std::string &containerId, ContainerData &results) {

    const std::string memStatFile = "/sys/fs/cgroup/cpu,cpuacct/system.slice/docker-" + containerId + ".scope/cpuacct.stat";

    std::ifstream infile(memStatFile);
    std::string name, value;
    std::map<std::string, long> dataMap;
    while (infile >> name >> value) {
        dataMap.insert({name, std::stol(value)});
        results.putMetricData("cpu",  name, std::stol(value));
    }
    results.putMetricData("cpu", "total", dataMap["user"] + dataMap["system"]);
}

#endif //DOCKER_AGENT_CPU_HPP
