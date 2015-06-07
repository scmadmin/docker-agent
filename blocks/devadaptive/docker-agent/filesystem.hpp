//
// Created by paul on 6/6/15.
//

#ifndef DOCKER_AGENT_FILESYSTEM_HPP
#define DOCKER_AGENT_FILESYSTEM_HPP

#include <fstream>
#include <string>


std::string getProcFile(const std::string& containerId) {
    return "/sys/fs/cgroup/devices/system.slice/docker-" + containerId + ".scope/cgroup.procs";
}

std::string parsePid(const std::string& procFile) {
    std::ifstream file(procFile);
    std::string str;
    if (std::getline(file, str)) {
        return str;
    } else {
        throw "couldn't read file";
    }
}

std::string getNamespacePath(const std::string& containerId) {
    const std::string& procFile = getProcFile(containerId);
    //todo: check it for validity
    const std::string& pid = parsePid(procFile);
    return "/proc/" + pid + "/ns/net";
}

#endif //DOCKER_AGENT_FILESYSTEM_HPP
