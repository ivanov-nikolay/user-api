### сервис, предоставляющий API для работы с данными пользователей

#### API
1. GET /user/{USER_ID) - Метод получения пользователя по идентификатору
2. POST /user - Метод добавления пользователя
3. DELETE /user/{USER_ID} - Метод удаления пользователя
4. PUT /user - Метод редактирования пользователя
5. GET /users - Метод поиска пользователей 
<br>
Может принимать query параметры:
<br>
Gender - строка male/female
<br>
Status - строка active/banned/deleted
<br>
FullName - строка 
<br>
SortAsk - true сортировка по возрастанию
<br>
SortDesc - true сортировка по убыванию 
<br>
AttributesToSort - строка - атрибут по которому сортировка id/name/surname/patronymic/gender/birthday/join_date
<br>
Limit - целое число
<br>
Offset - целое число
