<?php

namespace App;
/**
 * Class Handler
 * @package App
 */
class Handler
{
    /**
     * @param $data
     * @return
     */
    public function handle($input) {
 
        $catWords=array('tabby','lynx');
        $feedback = "Analysis NOT added to cat topic...";

        $headers = array('Accept' => 'application/json');
        $response = \Unirest\Request::post('http://gateway.openfaas:8080/function/inception', $headers, $input);
        
        $message = array("url" => $input
                        , "inferences" => json_decode(json_encode($response->body),true)
                        , "category" => "not-cats"
                        );
        
        foreach ($message['inferences'] as $inference){
            
            if ($inference['score'] < 0.2){
                break;
            }

            if (strpos($inference['name'], 'cat') !== false | in_array($inference['name'], $catWords)) {
              $feedback = "Analysis added to cat topic!";
              $message['category'] = "cats";
              break;
            }        
        }

        $config = \Kafka\ProducerConfig::getInstance();
        $config->setMetadataRefreshIntervalMs(10000);
        $config->setMetadataBrokerList('kafka.openfaas:9092');
        $config->setBrokerVersion('1.0.0');
        $config->setRequiredAck(1);
        $config->setIsAsyn(false);
        $config->setProduceInterval(500);
        $producer = new \Kafka\Producer();

        $producer->send([
            [
                'topic' => $message['category'],
                'value' => json_encode($message),
                'key' => '',
            ],
        ]);

        return $feedback;
    }
}
