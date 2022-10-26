# 解析Formats传入的MediaType
[{text html .}]
# Site实例
# 拥有了语言相关的配置项
# 每个页面类型也明确了自己的输出格式
Site:
&main.Site{
  language:(*main.Language)(0xc0000161b0),
  outputFormats:map[string]main.Formats{
    "404":main.Formats{
      main.Format{
        Name:"HTML",
        MediaType:main.Type{
          MainType:"text",
          SubType:"html",
          Delimiter:"."}, BaseName:"index"}},
    "page":main.Formats{
      main.Format{
        Name:"HTML",
        MediaType:main.Type{
          MainType:"text",
          SubType:"html",
          Delimiter:"."}, BaseName:"index"}}},
  outputFormatsConfig:main.Formats{
    main.Format{
      Name:"HTML",
      MediaType:main.Type{
        MainType:"text",
        SubType:"html",
        Delimiter:"."}, BaseName:"index"}},
  mediaTypesConfig:main.Types{
    main.Type{MainType:"text",
    SubType:"html", Delimiter:"."}}}

Program exited.