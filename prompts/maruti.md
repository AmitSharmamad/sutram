create tables to store chatbot conversations as mentioned below

1. ai_chat -> id, title, created_at
2. ai_message -> id, chat_id, query, response, feedback (like, dislike), feedback_comment
3. ai_actions -> id, message_id, type, value

a chat contains one or more messages
a message belongs to only one chat
a message contains zero or more ai_actions
an action belongs to only one message

create only sql, don't make any services or apis for this