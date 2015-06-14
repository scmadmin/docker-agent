#include "lasote/docker_client/client.h"
#include "filesystem.hpp"
#include "network.hpp"
#include "memory.hpp"
#include "container.hpp"
#include "ContainerData.hpp"
#include "HTTPClient.h"
#include "curl_client.hpp"


//std::string METRICS_END_POINT = "http://192.168.0.22:3000/porter/metrics";
std::string METRICS_END_POINT = "https://devadaptive.com:/p/porter/metrics";

ERR_F error_cb = [] (int status, string desc) {
    cout << "Error: " << status <<  endl  << desc;
};


void listContainers() {
    DockerClient client("http://localhost:4243");

    auto c6 = client.list_containers([] ( jsonxx::Object ret) {
        HTTPClient porterClient("http://192.168.0.22:3000/porter");
        JSON_F logResponse = [](jsonxx::Object ret) { cout << ret.json() << endl; };

        const std::string key = "data";
        const std::string tenantId = "tenant2";
        //cout << ret.json() << endl;
        JSON_ARRAY array = ret.get<JSON_ARRAY>(key);
        const std::vector<jsonxx::Value *> containers = array.values();

        for (auto &value : containers) {
            if (value->is<jsonxx::Object>()) {
                auto dockerContainer = value->get<jsonxx::Object>();
                if (dockerContainer.has<jsonxx::String>("Id")) {
                    auto id = dockerContainer.get<jsonxx::String>("Id");
                    ContainerData containerData{id, tenantId};
                    getContainerData(dockerContainer, containerData);
                    getNetworkData(id, containerData);
                    resetNSHack();
                    getMemoryData(id, containerData);
                    //cout << "metrics: " << containerData.getMetricArray().json() << endl;
                    std::string jsonOut(containerData.getMetricArray().json());
                    std::string out("[{\"the\" : \"quick\",\"brown\":\"fox\"},{\"jumped\": \"over\", \"the\" : \"lazy\"}]");
                    std::string responseBuffer;
                    postJSON(METRICS_END_POINT, jsonOut, responseBuffer);
                    //porterClient.postMetrics(out, logResponse, error_cb);
                    cout << "container id: "  << id << " response: " << responseBuffer << endl;
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