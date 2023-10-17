Z powodu braku czasu implementacja nie zawiera dwóch komponentów:
- nie zwraca najczęstszych słów a pełną, nieposortowaną mapę "słowo" => "licznik" dla danego URLa
- nie zaimplementowałem wyciągnięcia samych słów z kodu HTML (np. przy użyciu "html/tokenizer"), parsuję cały dokument

Kod został zaimplementowany z myślą o wywołaniu wewnątrz skrapera dodawania kolejnych URLi do kolejki i tym samym możlwiości zamienienia scrapera na crawlera.

# Uruchamianie

`go run -race ./cmd/`