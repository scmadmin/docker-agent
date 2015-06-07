//
// Created by paul on 6/6/15.
//

#ifndef DOCKER_AGENT_NETWORK_HPP
#define DOCKER_AGENT_NETWORK_HPP

#include <stdlib.h>
#include <sched.h>
#include <fcntl.h>
#include <unistd.h>
#include <iostream>
#include <fstream>
#include <string>

#include <sstream>
#include <vector>

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

void parseNetworkInNamespace(const std::string& nsFile) {

    int fileDescriptor = open(nsFile.c_str(), O_RDONLY);

    //this puts us inside the containers namespace.
    int setnsResult = setns(fileDescriptor, CLONE_NEWNET);

    string line;

    //this is the /proc/net/dev of the container!
    ifstream myfile ("/proc/net/dev");
    if (myfile.is_open())
    {
        while ( getline (myfile,line) )
        {
            if (line.find("eth0") != std::string::npos) {
                auto vec = xxx(line, ' ');
                cout << vec.size() << endl;
           }
        }
        myfile.close();
    }

    else cout << "Unable to open file";
    close(fileDescriptor);

}


#endif //DOCKER_AGENT_NETWORK_HPP
