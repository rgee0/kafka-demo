import sys
import json
import os
import tweepy
import requests
from tweetbuilder import TweetBuilder


def twitter_api():
    with open("/var/openfaas/secrets/consumer-key", "r") as consumerkey:
        consumer_key = consumerkey.read().strip()
    with open("/var/openfaas/secrets/consumer-secret", "r") as consumersecret:
        consumer_secret = consumersecret.read().strip()

    auth = tweepy.OAuthHandler(consumer_key, consumer_secret)

    with open("/var/openfaas/secrets/access-token", "r") as accesstoken:
        access_token = accesstoken.read().strip()
    with open("/var/openfaas/secrets/access-secret", "r") as accesssecret:
        access_secret = accesssecret.read().strip()

    auth.set_access_token(access_token, access_secret)
    api = tweepy.API(auth)

    return api


def tweet_image(url, message):
    api = twitter_api()
    filename = '/tmp/temp.jpg'
    request = requests.get(url, stream=True)
    if request.status_code == 200:
        with open(filename, 'wb') as image:
            for chunk in request:
                image.write(chunk)

        api.update_with_media(filename, status=message)
        os.remove(filename)
    else:
        print "Unable to download image"

def handle(req):
    try:
        event = json.loads(req)
    except ValueError:
        sys.exit("No JSON object received")
        return

    topName=event['inferences'][0]['name']
    topScore=int(event['inferences'][0]['score'] * 100)
    tweetcontent = TweetBuilder()
    tweet = "#notCatFacts " + tweetcontent.generate(topName,topScore)
    return tweet_image(event['url'], tweet)
