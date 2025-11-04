# Отношения 
## user
Хранит данные о пользователе (почта, пароль, аватарка, дата создания и обновления аккаунта).
### Зависимости:
{id} -> {username, email, password, avatar_id, created_at, updated_at} <br>
{avatar_id} -> {username, email, password, created_at, updated_at} 

## board
Хранит данные доски: пользователь-создатель, название, обложка, статус (закрытая или нет), даты создания и обновления, настройки доступа (публичная, по ссылке, приватная).
### Зависимости:
{id} -> {owner_user_id, title, image_id, archived, created_at, updated_at, visibility}

## board_member
Хранит данные о связи доски и пользователей, являющихся ее участниками: id пользователя и доски, роль участника (админ, участник, наблюдатель) и даты создания и обновления. 
### Зависимости:
{id} -> {user_id, board_id, role, created_at, updated_at} <br>
{user_id, board_id} -> {id, role, created_at, updated_at}

## list
Хранит данные о списке карточек: id доски, название, позицию по счету на доске, даты создания и обновления.
### Зависимости:
{id} -> {board_id, title, position, created_at, updated_at} <br>
{board_id, position} -> {id, title, created_at, updated_at}

## card
Хранит данные о карточке с заданием: id создателя карточки - участника доски, id списка, текст на карточке, позицию в списке по счету, даты создания и обновления, дату дедлайна.
### Зависимости:
{id} -> {author_board_member_id, list_id, content, position, created_at, updated_at, complete_before} <br>
{list_id, position} -> {id, author_board_member_id, content, created_at, updated_at, complete_before} 

## card_member
Хранит данные о связи между карточкой с задачей и участниками доски, занимающейся этой задачей(прикрепленных к ней): id карточки и участника доски, даты создания и обновления.
### Зависимости:
{id} -> {card_id, board_member_id, created_at, updated_at} <br>
{card_id, board_member_id} -> {id, created_at, updated_at}

## comment
Хранит данные о комментарии к карточке: id карточки, id создателя комментария - участника доски, текст и даты создания и обновления.
### Зависимости:
{id} -> {card_id, board_member_owner_id, content, created_at, updated_at}

## attachment
Хранит данные о вложении к карточке: id карточки, название, ссылку на вложение, позицию среди вложений на карточке, даты создания и обновления.
### Зависимости:
{id} -> {card_id, title, file_id, position, created_at, updated_at}

## checklist
Хранит данные о чеклисте на карточке: id карточки, название, даты создания и обновления. Только один чеклист может существовать у карточки.
### Зависимости:
{id} -> {card_id, title, created_at, updated_at} <br>
{card_id} -> {id, title, created_at, updated_at}

## checklist_point
Хранит данные о пункте чеклиста: id чеклиста, текст, статус (отмечен или нет), позицию по счету в чеклисте, даты создания и обновления.
### Зависимости:
{id} -> {checklist_id, content, checked, position, created_at, updated_at} <br>
{checklist_id, position} -> {id, content, checked, created_at, updated_at}

## upload
Хранит данные загруженных файлов: id файла, название, путь до файла, даты создания и обновления.
### Зависимости:
{id} -> {title, url, created_at, updated_at}

# Проверка на НФБК:
1-ая НФ: Все атрибуты являются атомарными<br>

2-ая НФ: Выполняется НФ1 и все неключевые элементы зависят от простого первичного ключа {id} <br>

3-я НФ: Выполняется НФ2 и отсутствуют зависимости между неключевыми атрибутами. Каждый атрибут зависит <b>только </b> от ключа.<br>

НФ Бойса-Кодда: Все детерминанты функциональных зависимостей отношений являются потенциальными ключами. Т.е. по: {avatar_id}; {user_id, board_id}; {board_id, position}, {list_id, position}; {card_id, board_member_id}; {checklist_id, position} можно однозначно определить атрибуты справа в соотвествующих им функциональных зависимостях. 