/**
 * Print a simple "Hello world!"
 *
 * @file main.cpp
 * @section LICENSE

    This code is under MIT License, http://opensource.org/licenses/MIT
 */

#include <iostream>

#include "lasote/docker_client/client.h"

#include "filesystem.hpp"
#include "network.hpp"

int main() {

    DockerClient client("http://localhost:4243");

    // Error callback for all examples
    ERR_F error_cb = [] (int status, string desc) {
        cout << "Error: " << status <<  endl  << desc;
    };

    auto c7 = client.list_containers([] ( jsonxx::Object ret) {
        const std::string key = "data";
        JSON_ARRAY array = ret.get<JSON_ARRAY>(key);
        const std::vector<jsonxx::Value*> containers =  array.values();

        for (auto &value : containers) {
            if (value->is<jsonxx::Object>()) {
                auto dockerContainer = value->get<jsonxx::Object>();
                if (dockerContainer.has<jsonxx::String>("Id")) {
                    auto id = dockerContainer.get<jsonxx::String>("Id");
                    cout << "id: " << id << endl;
                    parseNetworkInNamespace(getNamespacePath(id));
                }
            }
        }


        //cout << "Images: " << ret.json() << endl;
    }, error_cb);

    run_loop();
}
