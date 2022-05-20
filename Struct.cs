struct Option
{
    public string value, name, remarks;
    public int id;

    public Option(string value, string name, string remarks, int id)
    {
        this.value = value;
        this.name = name;
        this.remarks = remarks;
        this.id = id;
    }
}

struct QL
{
    public int id;
    public string url, token, value;

    public QL(int id,string url,string token,string value)
    {
        this.id=id;
        this.url=url;
        this.token=token;
        this.value=value;
    }
}