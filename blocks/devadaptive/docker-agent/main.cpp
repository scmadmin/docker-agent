#include "lasote/docker_client/client.h"
#include "filesystem.hpp"
#include "network.hpp"
#include "memory.hpp"
#include "container.hpp"
#include "ContainerData.hpp"
#include "HTTPClient.h"

std::string METRICS_END_POINT = "http://192.168.0.22:3000/porter";
//std::string METRICS_END_POINT = "https://devadaptive.com:/p/porter";

ERR_F error_cb = [] (int status, string desc) {
    cout << "Error: " << status <<  endl  << desc;
};


void listContainers() {
    DockerClient client("http://localhost:4243");

    auto c6 = client.list_containers([] ( jsonxx::Object ret) {
        HTTPClient porterClient(METRICS_END_POINT);
        JSON_F logResponse = [](jsonxx::Object ret) { cout << ret.json() << endl; };

        const std::string key = "data";
        const std::string tenantId = "456";
        //cout << ret.json() << endl;
        JSON_ARRAY array = ret.get<JSON_ARRAY>(key);
        const std::vector<jsonxx::Value *> containers = array.values();

        for (auto &value : containers) {
            if (value->is<jsonxx::Object>()) {
                auto dockerContainer = value->get<jsonxx::Object>();
                if (dockerContainer.has<jsonxx::String>("Id")) {
                    auto id = dockerContainer.get<jsonxx::String>("Id");
                    //cout << "container id: "  << id << endl;
                    ContainerData containerData{id, tenantId};
                    getContainerData(dockerContainer, containerData);
                    getNetworkData(id, containerData);
                    resetNSHack();
                    getMemoryData(id, containerData);
                    //cout << "metrics: " << containerData.getMetricArray().json() << endl;
                    porterClient.postMetrics(containerData.getMetricArray(), logResponse, error_cb);
                }
            }
        }

    }, error_cb);
}

int main() {
    uv_timer_t timer;
    uv_timer_init(uv_default_loop(), &timer);
    uv_timer_start(&timer, (uv_timer_cb) &listContainers, 1000, 15000);

    uv_run(uv_default_loop(), UV_RUN_DEFAULT);

}