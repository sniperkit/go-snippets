---
#layout: default
layout: page
title: "Life matters"
tagline: "Miscellaneous notes, about life, living, etc."
category : notes
tags : [draft]
published: true
maths: true
---

### Prêt immobilier

#### Calcul des mensualités (équation)

Soit :

* __m__ : mensualité
* K : capital emprunté
* t : taux annuel proportionnel
* n : nombre de mensualités

Le montant des mensualités peut être obtenu au moyen de l'équation suivante :

    m = [(K*t)/12] / [1-(1+(t/12))^-n]

{% raw %}
$$
\begin{align*}
m = \frac{\frac{K \times t}{12}}{1 - (1 + \frac{t}{12})^{-n}}
\end{align*}
$$
{% endraw %}

Exemple avec `bc` :

```
~$ bc -l
scale=2

k=300000; t=4.3/100; n=12*30; m = (k*t / 12) / (1 - (1 + t/12)^-n); m
1484.61
```

Calcul du cout total, et du cout du crédit :

```
cout_total = m*n ; cout_total
534461.15

cout_credit = cout_total - k ; cout_credit
234461.15
```

Estimation du montant potentiellement “épargnable” pour un loyer équivalent de
1100€ :

```
gain_si_loyer = (m-1100)*n ; gain_si_loyer
138461.15
```

Sur une seule ligne :

```
k=300000; t=4.3/100; n=12*30; m=(k*t/12)/(1-(1+t/12)^-n); cout_total=m*n; cout_credit=cout_total-k; gain_si_loyer = (m-1100)*n;  m ; cout_total ; cout_credit ; gain_si_loyer
  1484.61
534461.15
234461.15
138461.15
```

Trouvé équation @ [droit-finances.commentcamarche.com](http://droit-finances.commentcamarche.com/forum/affich-4149900-comment-calculer-mensualites-d-un-pret-immobilier#1).

Cf simulateur [linternaute.com/pratique/argent/immobilier/calculatrice/vos-mensualites/](http://www.linternaute.com/pratique/argent/immobilier/calculatrice/vos-mensualites/)

__TODO:__ Prendre en compte les frais de notaire, le coût de l'assurance, le coût
d'une éventuelle revente puis rachat d'un autre bien; taxes?.
__TODO:__ Lien vers cette vidéo Youtube là...

__EOF__
