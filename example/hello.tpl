{% for person in persons %}{% if not forloop.First %}
{% endif %}Hello, {{ person.name }}!{% endfor %}
