# Računanje determinante matrice primjenom Laplasovog razvoja po prvoj vrsti

# Članovi tima

- Mihajlo Kušljić, SW53-2016

# Osnovni pojmovi
- Matrica formata mn nad poljem F je funkcija Mmn koja preslikava skup uređenih parova {(𝑖,𝑗)|𝑖 ∈ {1,2,…,𝑚} ∧ 𝑗 ∈ {1,2,…,𝑛}} u skup F.
- Matrice formata nn, tj. matrice koje imaju isti broj redova i kolona, nazivamo kvadratne matrice reda n.
- Kvadratna podmatrica reda r matrice Mmn je kvadratna matrica reda r koja se dobija kada se iz matrice Mmn izbaci proizvoljnih (m - r) vrsta i (n - r) kolona.
- Determinanta je funkcija koja preslikava skup svih kvadratnih matrica u skup F, gdje je F polje nad kojim su definisane matrice.
- Minor reda r neke matrice Mmn je determinanta neke njene kvadratne podmatrice reda r.
- Jedna od metoda izračunavanja determinante matrice je Laplasov razvoj po vrsti/koloni. Ako razložimo matricu A po vrsi i, tada je: det(𝐴)= Σ[ (−1)^(i + j) * Aij * Mij ], gdje je Mij minor, tj. determinanta matrice koja se dobija kada se iz matrice A ukloni i-ta vrsta i j-ta kolona.

# Motivacija problema
Determinante imaju široku primjenu u mnogim oblastima matematike. U algebri determinante se koriste za opisivanje invertibilnih matrica (imaju inverznu matricu) i da se opiše rješenje sistema linearnih jednačina pomoću Kramerovog pravila. Determinante se koriste i da se izračunaju zapremine u vektorskoj analizi: apsolutna vrijednost determinante realnih vektora jednaka je zapremini paralelopipeda koji grade ti vektori. Takođe na osnovu determinante može se utvrditi linearna zavisnost realnih vektora. U matematičkoj analizi, pri računanju integrala po podskupovima Euklidskog prostora R^n, koriste se tzv. jakobijani koji predstavljaju determinante Jakobijevih matrica...

# Opis problema
Upotrebom programskih jezika Python i Go obezbjediti serijsku i paralelnu implementaciju računanja determinante kvadratne matrice reda n, primjenom Laplasovog razvoja po prvoj vrsti. Matrica se učitava iz teksutalnog fajla koji prati slijedeći format:

- Prvi red u fajlu sadrži red matrice (n)
- Narednih n redova u fajlu sadrže vrste matrice, gdje su elemeni razdvojeni barem jednim razmakom

Po učitavanju matrice izvršava se serijski i paralelni proračun determinante matrice. Rezultati se dodaju u CSV fajl koji ima slijedeće zaglavlje:

- n - red matrice čija je determinanta izračunata
- exec_time_ms - vrijeme izvršavanja proračuna u milisekundama
- serial - da li je proračun serijski, moguće vrijednosti su: true (serijski propračun), false (paralelizovan proračun)
- implementation - jezik u kojem je propračun implementiran, moguće vrijednosti su: python, go

Argumenti programa su putanje do (validnih) datoteka koje sadrže matrice čije determinante treba izračunati i putanja do CSV fajla u kojem se čuvaju rezultati. Pomoću programskog jezika Pharo Smalltalk omogućiti učitavanje generisanih rezultata iz CSV fajla i poređenje performansi implementacija u Python i Go jeziku na grafikonu (prikaz vremena izvršavanja u odnosu na red matrice).

# Koncept rješenja
Za računanje determinante kompletne matrice, kao i svakog od minora kreira se po jedan zadatak. Svi zadaci pristupaju istoj instanci matrice. Svaki zadatak računa determinantu svoje podmatrice primjenom Laplasovog razvoja po prvoj vrsti: za svaki element Aij u prvoj vrsti podmatrice čiju determinantu računa će pokrenuti novi zadatak koji računa odgovarajući minor Mij (pogledati formulu u prvom poglavlju). Kada se njegovi podzadaci završe, na osnovu vrijednosti minora koje su podzadaci generisali, sračunava vrijednost determinante za svoju podmatricu i tu vrijednost upisuje u varijablu u kojoj njegov roditeljski zadatak očekuje pojavu međurezultata. Na taj način se formira hijerarhija zadataka - determinanta matrice se računa rekurzivnom primjenom Laplasovog razvoja po prvoj vrsti, gdje je za jedan razvoj zadužen jedan zadatak. Svakom zadatku se prslijeđuje:

- matrix - polazna matrica
- begin_row_index - indeks koji u polaznoj matrici odgovara prvom redu podmatrice dodjeljene zadatku,
- column_indexes - indeksi koji u polaznoj matrici odgovaraju kolonama podmatrice dodjeljene zadatku,
- result_holder - promjenljiva u koju treba upisati rezultat kako bi mogao biti iskorišćen u roditeljskom zadatku.

Baza rekurzije je zadatak koji dobije kvadratnu podmatricu reda 1 (ima samo jednu kolonu). U tom slučaju vrijednost u matrici se proslijeđuje kao rezultat i ne dolazi do kreiranja podzadataka. Takođe moguće je uvesti konstantu CUTOFF koja predstavlja minimalan red matrice čija se determinanta računa paralelno. Za minore čiji je red manji od CUTOFF, determinanta se računa serijski (na isti način, rekurzivno). Ovo se dešava zato što je za male vrijednosti n trošak kreiranja i upravljanja zadacima veći od ubrzanja koje donosi podjela na zadatke.
Napomena: za implementaciju zadataka u programskog jeziku Python koristiće se standardna multiprocessing biblioteka, a u programskom jeziku Go koristiće se ugrađene go rutine. 