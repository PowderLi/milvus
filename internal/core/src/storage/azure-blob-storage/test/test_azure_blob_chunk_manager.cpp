#include "../AzureBlobChunkManager.h"

void print(std::string msg) {
}

int main() {
    azure::AzureBlobChunkManager::InitLog("info", print);
//    std::shared_ptr<azure::AzureBlobChunkManager> client_;
//    client_ = std::make_shared<azure::AzureBlobChunkManager>(
//            storage_config.access_key_id,
//            storage_config.access_key_value,
//            storage_config.address,
//            storage_config.useIAM);
}