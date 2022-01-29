from datetime import datetime
from flask import request
from flask import Flask
import maxminddb

app = Flask(__name__)

@app.route("/health", methods=['GET'])
def health():
    return {'now': datetime.now().isoformat()}


@app.route("/", methods=['POST'])
def hello():
    ip = request.get_json()['ip']
    countries = request.get_json()['contrycodes']
    with maxminddb.open_database('GeoLite2-Country.mmdb') as reader:
        matches = reader.get('152.216.7.110')
        for country in countries:             
            if matches:
                for matchCountry in matches['country']['names'].values():
                    if country == matchCountry:
                        return {'isInCountryList': True}
    return {'isInCountryList': False}
            

if __name__ == "__main__":
    app.run()