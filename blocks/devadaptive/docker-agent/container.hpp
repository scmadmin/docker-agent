//
// Created by paul on 6/7/15.
//

#ifndef DOCKER_AGENT_CONTAINER_HPP
#define DOCKER_AGENT_CONTAINER_HPP

#include <string>
#include <map>
#include <sstream>
#include "hjiang/jsonxx/jsonxx.h"
#include "ContainerData.hpp"

bool getContainerData(jsonxx::Object dockerContainer, ContainerData &containerData) {


    //container command
    if (dockerContainer.has<jsonxx::String>("Command")) {
        containerData.putDocumentData("container/Command",
                                      dockerContainer.get<jsonxx::String>("Command"));
    }
    //container create time
    if (dockerContainer.has<jsonxx::Number>("Created")) {
        jsonxx::Number value = dockerContainer.get<jsonxx::Number>("Created");
        containerData.putMetricData("container/Created", static_cast<long>(value));

    }
    //container Image
    if (dockerContainer.has<jsonxx::String>("Image")) {
        containerData.putDocumentData(
                "container/Image",
                dockerContainer.get<jsonxx::String>("Image")
        );

    }
    //container Names
    if (dockerContainer.has<jsonxx::Array>("Names")) {
        jsonxx::Array names = dockerContainer.get<jsonxx::Array>("Names");
        auto values = names.values();
        std::stringstream appender;
        appender << "[ ";
        for (auto &value : values) {
            if (value->is<jsonxx::String>()) {
                appender << value->get<jsonxx::String>() << ", ";
            }
        }
        appender << " ]";
        containerData.putDocumentData(
                "container/Names",
                appender.str()
        );
    }
    //container Status
    if (dockerContainer.has<jsonxx::String>("Status")) {
        containerData.putDocumentData("container/Status", dockerContainer.get<jsonxx::String>("Status"));

    }
    return true;
}

#endif //DOCKER_AGENT_CONTAINER_HPP
