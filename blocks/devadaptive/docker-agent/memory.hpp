//
// Created by paul on 6/7/15.
//

#ifndef DOCKER_AGENT_MEMORY_HPP
#define DOCKER_AGENT_MEMORY_HPP

#include <string>
#include <map>
#include <fstream>


void getMemoryData(std::string &containerId, ContainerData &results) {

    const std::string memStatFile = "/sys/fs/cgroup/memory/system.slice/docker-" + containerId + ".scope/memory.stat";

    std::ifstream infile(memStatFile);
    std::string name, value;
    while (infile >> name >> value) {
        results.putMetricData("memory",  name, std::stol(value));
    }
}

#endif //DOCKER_AGENT_MEMORY_HPP
