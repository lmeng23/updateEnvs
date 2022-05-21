use reqwest::{get, header::HeaderMap, Client};
use serde_json::{Map, Value};
use std::collections::HashMap;
use std::time::Duration;

#[tokio::main]
async fn main() {
    let source_url = "";
    let source_client_id = "";
    let source_client_secret = "";

    let dist_url = "";
    let dist_client_id = "";
    let dist_client_secret = "";

    let source_token = get_token(source_url, source_client_id, source_client_secret)
        .await
        .replace('"', "");

    let cookie_info = get_envs(source_url, source_token.as_str()).await;
    let cookie_array = cookie_info["data"].as_array().unwrap();
    let mut cookies: Vec<String> = vec![];
    for item in cookie_array {
        cookies.push(item["value"].to_string().replace('"', ""));
    }
    let cookie_str = cookies.join("&");

    let dist_token = get_token(dist_url, dist_client_id, dist_client_secret)
        .await
        .replace('"', "");
    let dist_cookies_info = get_envs(dist_url, dist_token.as_str()).await;
    let dist_cookies_array = dist_cookies_info["data"].as_array().unwrap();

    let mut id = String::new();
    for item in dist_cookies_array {
        if item["name"].as_str().unwrap() == "JD_COOKIE" {
            id = item["id"].to_string();
        }
    }

    let res = update_envs(
        id.as_str(),
        dist_url,
        dist_token.as_str(),
        cookie_str.as_str(),
    )
    .await;

    if res["code"].to_string().contains("200") {
        println!("更新JD_COOKIE环境变量成功!");
    } else {
        println!("更新JD_COOKIE环境变量失败!");
    }
}

async fn get_token(url: &str, client_id: &str, client_secret: &str) -> String {
    let res = get(format!(
        "{}/open/auth/token?client_id={}&client_secret={}",
        url, client_id, client_secret
    ))
    .await
    .unwrap()
    .json::<Map<String, Value>>()
    .await
    .unwrap();

    format!("{} {}", res["data"]["token_type"], res["data"]["token"])
}

async fn get_envs(url: &str, token: &str) -> Map<String, Value> {
    let mut headers = HeaderMap::new();
    headers.insert("Accept", "application/json".parse().unwrap());
    headers.insert("Authorization", token.parse().unwrap());

    let client = Client::new();
    client
        .get(format!("{}/open/envs", url))
        .headers(headers)
        .timeout(Duration::from_secs(3))
        .send()
        .await
        .unwrap()
        .json::<Map<String, Value>>()
        .await
        .unwrap()
}

async fn update_envs(id: &str, url: &str, token: &str, value: &str) -> Map<String, Value> {
    let mut headers = HeaderMap::new();
    headers.insert("Accept", "application/json".parse().unwrap());
    headers.insert("Content-Type", "application/json".parse().unwrap());
    headers.insert("Authorization", token.parse().unwrap());

    let mut data = HashMap::new();
    data.insert("value", value);
    data.insert("name", "JD_COOKIE");
    data.insert("remarks", "京东Cookie");
    data.insert("id", id);

    let client = Client::new();
    client
        .put(format!("{}/open/envs", url))
        .headers(headers)
        .json(&data)
        .send()
        .await
        .unwrap()
        .json::<Map<String, Value>>()
        .await
        .unwrap()
}
