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
        typedef std::map<std::string, std::string> ConfigInfo;
        ConfigInfo configValues;
        std::string line;
        while (std::getline(infile, line)) {
            std::istringstream lineStream(line);
            std::string key;
            if (std::getline(lineStream, key, '=')) {
                std::string value;
                if (key[0] == '#')
                    continue;
                if (std::getline(lineStream, value)) {
                    properties[key] = value;
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
    std::map<std::string, std::string> properties;


};


#endif //DOCKER_AGENT_CONFIGURATION_H
