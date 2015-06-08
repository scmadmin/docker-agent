//
// Created by paul on 6/6/15.
//

#ifndef DOCKER_AGENT_FILESYSTEM_HPP
#define DOCKER_AGENT_FILESYSTEM_HPP

#include <fstream>
#include <string>
#include <stdlib.h>
#include <sched.h>
#include <fcntl.h>
#include <unistd.h>

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

int openNamespace(const string &containerId, int &fileDescriptor) {
    auto nsFile = getNamespacePath(containerId);
    fileDescriptor = open(nsFile.c_str(), O_RDONLY);
    if (-1 == fileDescriptor) {
        throw "unable to open file " + nsFile;
    }
    //this puts us inside the containers namespace.
    int nsRet = setns(fileDescriptor, CLONE_NEWNET);
    if (-1 == nsRet) {
        throw "unable to set namespace based on " + nsFile;
    }
    //cout << "opening " << fileDescriptor << endl;
}

void closeNamespace(int fileDescriptor) {
    close(fileDescriptor);
}

#endif //DOCKER_AGENT_FILESYSTEM_HPP
