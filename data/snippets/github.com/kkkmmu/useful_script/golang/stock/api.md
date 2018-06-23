和讯网股票接口：
日k（所有数据） 
http://flashquote.stock.hexun.com/Quotejs/DA/1_601398_DA.html

指定时间区域的日k 
http://webstock.quote.hermes.hexun.com/a/kline?code=sse601398&start=20170909150000&number=-1000&type=5&callback=callback

分时线 
http://flashquote.stock.hexun.com/Quotejs/MA/1_000001_MA.html?

指定时间区域的分时线 
http://webstock.quote.hermes.hexun.com/a/minute?code=sse600000&start=20170424000000&number=6000&callback=callback

明细接口 
http://flashquote.stock.hexun.com/Stock_DL.ASPX?m=1&c=000001&t=0.6584310308098793






网易的接口是：
网易的数据格式为csv文件

日线

http://quotes.money.163.com/service/chddata.html?code=代码&start=开始时间&end=结束时间&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;TURNOVER;VOTURNOVER;VATURNOVER;TCAP;MCAP

参数说明：代码为股票代码，上海股票前加0，如600756变成0600756，深圳股票前加1。时间都是6位标识法，年月日，fields标识想要请求的数据。可以不变。

例如大盘指数数据查询（上证指数000001前加0，沪深300指数000300股票前加0，深证成指399001前加1，中小板指399005前加1，创业板指399006前加1）： 
http://quotes.money.163.com/service/chddata.html?code=0000300&start=20151219&end=20171108&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;VOTURNOVER

上海股票数据查询（浪潮）：http://quotes.money.163.com/service/chddata.html?code=0600756&start=20160902&end=20171108&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;VOTURNOVER;









股票代码的信息是在东方财富网中获取(http://quote.eastmoney.com/stocklist.html)


和讯获取个行业数据:
http://quote.stock.hexun.com/hqzx/industryquotestd.aspx?sorttype=4&page=1&count=50&time=00080


股票实时数据:http://hq.sinajs.cn/list=sh601006
