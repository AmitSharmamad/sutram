INSERT INTO sutram.user (
  id, name, username, email, phone_number, password, is_active
) VALUES (
  'e65c67f1-2b43-4581-b33b-ec26295e741b',
  'system',
  'system',
  'system@sutram.com',
  '-',
  'system',
  true
);

INSERT INTO sutram.strategy (
  id, user_id, name, description, entry_expression,
exit_expression
) VALUES (
  'c524747e-cb18-4f2e-86df-77b19dc6d23a',
  'e65c67f1-2b43-4581-b33b-ec26295e741b',
  'system strategy',
  'system strategy',
  '{"type":"ExpressionStatement","returns":"bool","value":{"type":"LogicalExpression","operator":"AND","returns":"bool","lhs":{"type":"IndicatorExpression","operator":"CrossesAbove","returns":"bool","lhs":{"type":"Indicator","name":"EMA","point":"close","returns":"float","value":{"period":5,"offset":0}},"rhs":{"type":"Number","name":"Number","point":"close","returns":"float","value":{"value":21000}}},"rhs":{"type":"IndicatorExpression","operator":"CrossesAbove","returns":"bool","lhs":{"type":"Indicator","name":"EMA","point":"close","returns":"float","value":{"value":{"period":5,"offset":0}}},"rhs":{"type":"Number","name":"Number","point":"close","returns":"float","value":{"value":21000}}}}}',
  '{"type":"ExpressionStatement","returns":"bool","value":{"type":"LogicalExpression","operator":"AND","returns":"bool","lhs":{"type":"IndicatorExpression","operator":"CrossesAbove","returns":"bool","lhs":{"type":"Indicator","name":"EMA","point":"close","returns":"float","value":{"period":5,"offset":0}},"rhs":{"type":"Number","name":"Number","point":"close","returns":"float","value":{"value":21000}}},"rhs":{"type":"IndicatorExpression","operator":"CrossesAbove","returns":"bool","lhs":{"type":"Indicator","name":"EMA","point":"close","returns":"float","value":{"value":{"period":5,"offset":0}}},"rhs":{"type":"Number","name":"Number","point":"close","returns":"float","value":{"value":21000}}}}}'
);
