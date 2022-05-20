using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using System.Collections;
using System.Net.Http.Headers;


string sourceUrl = "";
string sourceClientID = "";
string sourceClientSecret = "";

string distUrl = "";
string distClientID = "";
string distClientSecret = "";


string sourceToken = await GetTokenAsync(sourceUrl, "", "");
JArray? cookieArray = await GetEnvsAsync(sourceUrl, sourceToken);

ArrayList cookies = new ArrayList();
foreach (JObject item in cookieArray)
{
    cookies.Add(item["value"].ToString());
}
string cookieStr = string.Join("&", (string[])cookies.ToArray(typeof(string)));

string distToken = await GetTokenAsync(distUrl, "", "");
JArray? cookiesInfo = await GetEnvsAsync(distUrl, distToken);

int id = 0;
foreach (JObject item in cookiesInfo)
{
    if (item["name"].ToString() == "JD_COOKIE")
    {
        id = Convert.ToInt32(item["id"].ToString());
    }
}

JObject? res = await UpdateEnvsAsync(new QL(id, distUrl, distToken, cookieStr));

if (Convert.ToInt32(res["code"].ToString()) == 200)
{
    Console.WriteLine("更新JD_COOKIE环境变量成功!");
}
else
{
    Console.WriteLine("更新JD_COOKIE环境变量失败!");
}

Console.ReadLine();


static async Task<string> GetTokenAsync(string url, string client_id, string client_secret)
{
    url = $"{url}/open/auth/token?client_id={client_id}&client_secret={client_secret}";
    using var httpClient = new HttpClient();
    using var request = new HttpRequestMessage(new HttpMethod("Get"), url);
    var response = await httpClient.SendAsync(request);
    using HttpContent content = response.Content;
    JObject? res = JsonConvert.DeserializeObject(await content.ReadAsStringAsync()) as JObject;
    return $"{res["data"]["token_type"]} {res["data"]["token"]}";
}


static async Task<JArray?> GetEnvsAsync(string url, string token)
{
    url = $"{url}/open/envs";
    using var httpClient = new HttpClient();
    using var request = new HttpRequestMessage(new HttpMethod("Get"), url);
    request.Headers.TryAddWithoutValidation("Accept", "application/json");
    request.Headers.TryAddWithoutValidation("Authorization", token);
    var response = await httpClient.SendAsync(request);
    using HttpContent content = response.Content;
    JObject? res = JsonConvert.DeserializeObject(await content.ReadAsStringAsync()) as JObject;
    return res["data"] as JArray;
}


static async Task<JObject?> UpdateEnvsAsync(QL ql)
{
    string url = $"{ql.url}/open/envs";
    using var httpClient = new HttpClient();
    using var request = new HttpRequestMessage(new HttpMethod("PUT"), url);
    request.Headers.TryAddWithoutValidation("Accept", "application/json");
    request.Headers.TryAddWithoutValidation("Content-Type", "application/json");
    request.Headers.TryAddWithoutValidation("Authorization", ql.token);
    Option option = new(ql.value, "JD_COOKIE", "京东Cookie", ql.id);
    request.Content = new StringContent(JsonConvert.SerializeObject(option));
    request.Content.Headers.ContentType = MediaTypeHeaderValue.Parse("application/json");
    var response = await httpClient.SendAsync(request);
    using HttpContent content = response.Content;
    JObject? res = JsonConvert.DeserializeObject(await content.ReadAsStringAsync()) as JObject;
    return res;
}