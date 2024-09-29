import ipaddress
import maxminddb
from fastapi import FastAPI, Request

city_reader = maxminddb.open_database('GeoLite2-City.mmdb')
asn_reader = maxminddb.open_database('GeoLite2-ASN.mmdb')
cn_reader = maxminddb.open_database('GeoCN.mmdb')
lang = ["zh-CN","en"]
asn_map = {
    9812:"东方有线",
    9389:"中国长城",
    17962:"天威视讯",
    17429:"歌华有线",
    7497:"科技网",
    24139:"华数",
    9801:"中关村",
    4538:"教育网",
    24151:"CNNIC",
    
    38019:"中国移动",139080:"中国移动",9808:"中国移动",24400:"中国移动",134810:"中国移动",24547:"中国移动",
    56040:"中国移动",56041:"中国移动",56042:"中国移动",56044:"中国移动",132525:"中国移动",56046:"中国移动",
    56047:"中国移动",56048:"中国移动",59257:"中国移动",24444:"中国移动",
    24445:"中国移动",137872:"中国移动",9231:"中国移动",58453:"中国移动",
    
    4134:"中国电信",4812:"中国电信",23724:"中国电信",136188:"中国电信",137693:"中国电信",17638:"中国电信",
    140553:"中国电信",4847:"中国电信",140061:"中国电信",136195:"中国电信",17799:"中国电信",139018:"中国电信",
    133776:"中国电信",58772:"中国电信",146966:"中国电信",63527:"中国电信",58539:"中国电信",58540:"中国电信",
    141998:"中国电信",138169:"中国电信",139203:"中国电信",58563:"中国电信",137690:"中国电信",63838:"中国电信",
    137694:"中国电信",137698:"中国电信",136167:"中国电信",148969:"中国电信",134764:"中国电信",
    134770:"中国电信",148981:"中国电信",134774:"中国电信",136190:"中国电信",140647:"中国电信",
    132225:"中国电信",140485:"中国电信",4811:"中国电信",131285:"中国电信",137689:"中国电信",
    137692:"中国电信",140636:"中国电信",140638:"中国电信",140345:"中国电信",38283:"中国电信",
    140292:"中国电信",140903:"中国电信",17897:"中国电信",134762:"中国电信",139019:"中国电信",
    141739:"中国电信",141771:"中国电信",134419:"中国电信",140276:"中国电信",58542:"中国电信",
    140278:"中国电信",139767:"中国电信",137688:"中国电信",137691:"中国电信",4809:"中国电信",
    58466:"中国电信",137687:"中国电信",134756:"中国电信",134760:"中国电信",
    133774:"中国电信",133775:"中国电信",4816:"中国电信",134768:"中国电信",
    58461:"中国电信",58519:"中国电信",58520:"中国电信",131325:"中国电信",

    4837:"中国联通",4808:"中国联通",134542:"中国联通",134543:"中国联通",10099:"中国联通",
    140979:"中国联通",138421:"中国联通",17621:"中国联通",17622:"中国联通",17816:"中国联通",
    140726:"中国联通",17623:"中国联通",136958:"中国联通",9929:"中国联通",58519:"中国联通",
    140716:"中国联通",4847:"中国联通",136959:"中国联通",135061:"中国联通",139007:"中国联通",

    59019:"金山云",
    135377:"优刻云",
    45062:"网易云",
    137718:"火山引擎",
    37963:"阿里云",45102:"阿里云国际",
    45090:"腾讯云",132203:"腾讯云国际",
    55967:"百度云",38365:"百度云",
    58519:"华为云", 55990:"华为云",136907:"华为云",
    4609:"澳門電訊",
    134773:"珠江宽频",
    1659:"台湾教育网",
    8075:"微软云",
    17421:"中华电信",
    3462:"HiNet",
    13335:"Cloudflare",
    55960:"亚马逊云",14618:"亚马逊云",16509:"亚马逊云",
    15169:"谷歌云",396982:"谷歌云",36492:"谷歌云",
}

def get_as_info(number):
    r = asn_map.get(number)
    if r:
        return r
    
def get_des(d):
    for i in lang:
        if i in d['names']:
            return d['names'][i]
    return d['names']['en']

def get_country(d):
    r = get_des(d)
    if r in ["香港", "澳门", "台湾"]:
        return "中国" + r
    return r

def province_match(s):
    arr=['内蒙古','黑龙江','河北','山西','吉林','辽宁','江苏','浙江','安徽','福建','江西','山东','河南','湖北','湖南','广东','海南','四川','贵州','云南','陕西','甘肃','青海','广西','西藏','宁夏','新疆','北京','天津','上海','重庆']
    for i in arr:
        if i in s:
            return i
    return ''

def de_duplicate(regions):
    regions = filter(bool,regions)
    ret = []
    [ret.append(i) for i in regions if i not in ret]
    return ret

def get_addr(ip, mask):
    network = ipaddress.ip_network(f"{ip}/{mask}", strict=False)
    first_ip = network.network_address
    return f"{first_ip}/{mask}"

def get_maxmind(ip: str):
    ret = {"ip":ip}
    asn_info = asn_reader.get(ip)
    if asn_info:
        as_ = {"number":asn_info["autonomous_system_number"],"name":asn_info["autonomous_system_organization"]}
        info = get_as_info(as_["number"])
        if info:
            as_["info"] = info
        ret["as"] = as_

    city_info, prefix = city_reader.get_with_prefix_len(ip)
    ret["addr"] = get_addr(ip, prefix)
    if not city_info:
        return ret
    if "country" in city_info:
        country_code = city_info["country"]["iso_code"]
        country_name = get_country(city_info["country"])
        ret["country"] = {"code":country_code,"name":country_name}
    
    if "registered_country" in city_info:
        registered_country_code = city_info["registered_country"]["iso_code"]
        ret["registered_country"] = {"code":registered_country_code,"name":get_country(city_info["registered_country"])}
        
    regions = [get_des(i) for i in city_info.get('subdivisions', [])]

    
    if "city" in city_info:
        c = get_des(city_info["city"])
        if (not regions or c not in regions[-1])and c not in country_name:
            regions.append(c)
            
    regions = de_duplicate(regions)
    if regions:
        ret["regions"] = regions
    
    return ret

def get_cn(ip:str, info={}):
    ret, prefix = cn_reader.get_with_prefix_len(ip)
    if not ret:
        return
    info["addr"] = get_addr(ip, prefix)
    regions = de_duplicate([ret["province"],ret["city"],ret["districts"]])
    if regions:
        info["regions"] = regions
        info["regions_short"] = de_duplicate([province_match(ret["province"]),ret["city"].replace('市',''),ret["districts"]])
    if "as" not in info:
        info["as"] = {}
    info["as"]["info"] = ret['isp']
    if ret['net']:
        info["type"] = ret['net']
    return ret

def get_ip_info(ip):
    info = get_maxmind(ip)
    if "country" in info and info["country"]["code"] == "CN" and ("registered_country" not in info or info["registered_country"]["code"] == "CN"):
        get_cn(ip,info)
    return info

def query():
    while True:
        try:
            ip = input('IP：   \t').strip()
            info = get_ip_info(ip)
                
            print(f"网段：\t{info['addr']}")
                
            if "as" in info:
                print(f"ISP：\t",end=' ')
                if "info" in info["as"]:
                    print(info["as"]["info"],end=' ')
                else:
                    print(info["as"]["name"],end=' ')
                if "type" in info:
                    print(f"({info['type']})",end=' ')
                print(f"ASN{info['as']['number']}",end=' ')
                print(info['as']["name"])
                
            if "registered_country" in info and ("country" not in info or info["country"]["code"] != info["registered_country"]["code"]):
                print(f"注册地：\t{info['registered_country']['name']}")
                
            if "country" in info:
                print(f"使用地：\t{info['country']['name']}")
                
            if "regions" in info:
                print(f"位置：    \t{' '.join(info['regions'])}")
                
        except Exception as e:
            print(e)
            raise e
        finally:
            print("\n")
            
app = FastAPI()

@app.get("/")
def api(request: Request, ip: str = None):
    if not ip:
        ip = request.headers.get("x-forwarded-for") or request.headers.get("x-real-ip") or request.client.host
    return get_ip_info(ip.strip())

@app.get("/{ip}")
def path_api(ip):
    return get_ip_info(ip)

if __name__ == '__main__':
    query()
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080, server_header=False, proxy_headers=True)
