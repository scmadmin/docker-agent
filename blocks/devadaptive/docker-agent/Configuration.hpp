//
// Created by paul on 6/28/15.
//

#ifndef DOCKER_AGENT_CONFIGURATION_H
#define DOCKER_AGENT_CONFIGURATION_H


#include <stdlib.h>
#include <fstream>
#include <map>
#include "utils.hpp"

class Configuration {


public:


    Configuration(const std::string &propertyFile) {
        std::ifstream infile(propertyFile);
        std::string line;
        std::map<std::string, std::string> properties;
        while (getline(infile, line)) {
            //todo: trim start of line
            if (line.length() > 0 && line[0] != '#') {
                std::istringstream lineStream(line);
                std::string key, value;
                char delim;
                if ((lineStream >> key >> delim >> value) && (delim == '=')) {
                    std::cout << key << " = " << value << std::endl;
                    properties[key] = value;
                } else {
                    //std::cout << key << " = " << value << std::endl;
                }
            }
        }
    }

public:
    const std::string &getAgentId() const {
        return agentId;
    }

    const std::string &getTenantId() const {
        return tenantId;
    }

    const std::string &getUrl() const {
        return url;
    }

private:
    std::string agentId;
    std::string tenantId;
    std::string url = "https://devadaptive.com/p/porter/metrics";
    bool debug = false;


};


#endif //DOCKER_AGENT_CONFIGURATION_H
