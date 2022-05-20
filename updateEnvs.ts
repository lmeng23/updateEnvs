import axios from 'axios';

const sourceUrl = '';
const sourceClientID = '';
const sourceClientSecret = '';

const distUrl = '';
const distClientID = '';
const distClientSecret = '';

!(async () => {
    const sourceToken = await getToken(sourceUrl, sourceClientID, sourceClientSecret);
    const cookies: string[] = [];

    const cookieObject = await getEnvs(sourceUrl, sourceToken);
    for (const item of cookieObject as any) {
        cookies.push(item.value);
    }

    const cookieStr = cookies.join('&');

    const distToken = await getToken(distUrl, distClientID, distClientSecret);

    const cookiesInfo = await getEnvs(distUrl, distToken);

    let id: any;

    for (const item of cookiesInfo as any) {
        if (item.name === 'JD_COOKIE') {
            id = item.id;
        }
    }

    const res = await updateEnvs(id as number, distUrl, distToken, cookieStr);

    if ((res as any).code === 200) {
        console.log('更新JD_COOKIE环境变量成功!');
    } else {
        console.log('更新JD_COOKIE环境变量失败!');
    }


})()

// 获取青龙面板登录token
async function getToken(url: string, client_id: string, client_secret: string): Promise<string> {

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

async function getEnvs(url: string, token: string): Promise<object> {

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
async function updateEnvs(id: number, url: string, token: string, value: string): Promise<object> {
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