//
// Created by paul on 6/6/15.
//

#ifndef DOCKER_AGENT_NETWORK_HPP
#define DOCKER_AGENT_NETWORK_HPP


#include <iostream>
#include <fstream>
#include <string>

#include <sstream>
#include <vector>

#include "filesystem.hpp"
#include "ContainerData.hpp"


std::vector<std::string> &split2(const std::string &s, char delim, std::vector<std::string> &elems) {
    std::stringstream ss(s);
    std::string item;
    while (std::getline(ss, item, delim)) {
        if (!item.empty())
            elems.push_back(item);
    }
    return elems;
}

std::vector<std::string> xxx(const std::string &s, char delim) {
    std::vector<std::string> elems;
    split2(s, delim, elems);
    return elems;
}

void getNetworkData(const std::string &containerId, ContainerData& containerData) {

    int fileDescriptor;
    openNamespace(containerId, fileDescriptor);
    //cout << "opened fd = " << fileDescriptor << endl;

    string line;

    //this is the /proc/net/dev of the container!
    ifstream myfile("/proc/net/dev");
    if (myfile.is_open()) {
        while (getline(myfile, line)) {
            if (line.find("eth0") != std::string::npos) {
                auto vec = xxx(line, ' ');
                //cout << "Network: rx: " << vec[1] << " tx: " << vec[9] << endl;
                containerData.putMetricData("network/receive_size", std::stol(vec[1]));
                containerData.putMetricData("network/transmit_size", std::stol(vec[9]));
            }
        }
        myfile.close();
    }

    else cout << "Unable to open file";
    closeNamespace(fileDescriptor);

}




#endif //DOCKER_AGENT_NETWORK_HPP
