const axios = require('axios')

const sourceUrl = '';
const sourceClientID = '';
const sourceClientSecret = ''

const distUrl = '';
const distClientID = ''
const distClientSecret = ''

    !(async () => {
        const sourceToken = await getToken(sourceUrl, sourceClientID, sourceClientSecret);
        const cookies = [];

        cookieObject = await getEnvs(sourceUrl, sourceToken);
        cookieObject.forEach(item => {
            cookies.push(item.value);
        });

        const cookieStr = cookies.join('&');

        const distToken = await getToken(distUrl, distClientID, distClientSecret);

        const cookiesInfo = await getEnvs(distUrl, distToken);

        let id;
        cookiesInfo.forEach(item => {
            if (item.name === 'JD_COOKIE') {
                id = item.id;
            }
        });

        const res = await updateEnvs(id, distUrl, distToken, cookieStr);

        if (res.code === 200) {
            console.log('更新JD_COOKIE环境变量成功!');
        }


    })()

// 获取青龙面板登录token
async function getToken(url, client_id, client_secret) {

    const {
        data: res
    } = await axios.get(`${url}/open/auth/token`, {
        params: {
            "client_id": client_id,
            "client_secret": client_secret
        }
    });

    return `${res.data['token_type']} ${res.data['token']}`;
}

// 获取青龙面板JD_COOKIE环境变量
async function getEnvs(url, token) {

    const {
        data: res
    } = await axios.get(`${url}/open/envs`, {
        headers: {
            "Accept": "application/json",
            "Authorization": token
        }
    });

    return res.data;
}

// 更新青龙面板JD_COOKIE环境变量
async function updateEnvs(id, url, token, value) {
    const {
        data: res
    } = await axios.put(`${url}/open/envs`, {
        value: value,
        name: "JD_COOKIE",
        remarks: '京东Cookie',
        id: id
    }, {
        headers: {
            "Accept": "application/json",
            "Content-Type": "application/json",
            "Authorization": token
        }
    });
    return res;
}