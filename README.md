# Public_OpenAPI
OpenAPI를 다루고 기록으로 남겨두려고 합니다.
* * *

data_seoul 폴더
=============
서울 열린데이터광장 OpenAPI를 받아왔습니다.
https://data.seoul.go.kr/
* * *

garak_co_kr 폴더
=============
서울특별시 농수산 식품공사 OpenAPI를 받아왔습니다.
https://www.garak.co.kr/

### JSON 정보 문제

garak.co.kr에서 받은 JSON 정보의 key가 없는 문제와 배열의 형식이 이상한 문제가 있었습니다.
**ex**
{ LIST_COUNT:0000,:[{정보1}],[{정보2}],[{정보3}] }

해결과정은 아래와 같습니다.

LIST_COUNT가 쌍따옴표로 묶여있지 않은 문제
<pre>
<code>
newData:= strings.ReplaceAll(string(body), "{ LIST_COUNT", "{ \"LIST_COUNT\"")
</code>
</pre>

배열 정보의 key가 없던 문제
<pre>
<code>
newData = strings.ReplaceAll(newData, ",:", ",\"items\":")
</code>
</pre>

배열 속 정보가 대괄호로 묶여있는 문제
<pre>
<code>
newData = strings.ReplaceAll(newData, "],[", ",")
</code>
</pre>
* * *
