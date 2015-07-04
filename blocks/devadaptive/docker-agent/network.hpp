//
// Created by paul on 6/6/15.
//

#ifndef DOCKER_AGENT_NETWORK_HPP
#define DOCKER_AGENT_NETWORK_HPP

#include <pthread.h>
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


void *getNetworkData2X(void *blah) {
    cout << "opening nsfile " << endl;
    int fileDescriptor = open("/proc/53823/ns/net", O_RDONLY);
    cout << "file open " << fileDescriptor << endl;
    int setnsResult = setns(fileDescriptor, CLONE_NEWNET);
    cout << setnsResult << endl;
    string line;
    ifstream myfile("/proc/net/dev");
    cout << "proc net dev ctor" << endl;
    if (myfile.is_open()) {
        while (getline(myfile, line)) {
            cout << line << endl;
        }
        myfile.close();
        cout << "closing proc net dev " << endl;
    } else
        cout << "Unable to open file";
    close(fileDescriptor);
    cout << "closed ns file" << endl;
}

void resetNSHack() {
    //cout << "hack opening nsfile " << endl;
    int fileDescriptor = open("/proc/1/ns/net", O_RDONLY);
    //cout << "hack file open " << fileDescriptor << endl;
    int setnsResult = setns(fileDescriptor, CLONE_NEWNET);
    //cout << setnsResult << endl;
    close(fileDescriptor);
    //cout << "hack closed ns file" << endl;
}

void getNetworkData2() {
    pthread_t thread1;
    int blah = 0;
    getNetworkData2X(&blah);
    resetNSHack();
    //int r = pthread_create(&thread1, NULL, &getNetworkData2X, (void *) &blah);
    //pthread_join(thread1, NULL);
}

void getNetworkData(const std::string &containerId, ContainerData &containerData) {

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
                containerData.putMetricData("network", "receive_size", std::stol(vec[1]));
                containerData.putMetricData("network", "transmit_size", std::stol(vec[9]));
            }
        }
        myfile.close();
    } else {
        string error = "Unable to open file /proc/dev/net";
        cout << error << endl;
        throw error;
    }
    closeNamespace(fileDescriptor);

}


#endif //DOCKER_AGENT_NETWORK_HPP
