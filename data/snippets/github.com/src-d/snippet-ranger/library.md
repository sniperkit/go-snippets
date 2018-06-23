
# Topic modeling

## BigARTM
* [BigARTM](http://bigartm.org) main site
* BigARTM [Tutorial on Probabilistic Topic Modeling: Additive Regularization for Stochastic Matrix Factorization](http://www.machinelearning.ru/wiki/images/1/1f/Voron14aist.pdf)
* BigARTM [Additive Regularization for Topic Models of Text Collections](http://www.machinelearning.ru/wiki/images/2/21/Voron14dan-eng.pdf)
* Main idea: Additive Regularization of Topic Models
* hARTM for hierarchical topic modeling. [presentation](http://www.machinelearning.ru/wiki/images/d/dc/2.Chirkova.pdf)

## Other
* [Collaborative deep learning for recommender systems](https://arxiv.org/pdf/1409.2944.pdf). [code](https://github.com/js05212/MXNet-for-CDL/blob/master/collaborative-dl.ipynb).
> We generalize recent advances in deep learning from i.i.d. input to non-i.i.d. (CF-based) input and propose in this paper a hierarchical Bayesian model called collaborative deep learning (CDL), which jointly performs deep representation learning for the content information and collaborative filtering for the ratings (feedback) matrix. Extensive experiments on three real-world datasets from different domains show that CDL can significantly advance the state of the art.
* [Scalable Deep Poisson Factor Analysis for Topic Modeling](http://proceedings.mlr.press/v37/gan15.pdf), [matlab code](https://github.com/zhegan27/dpfa_icml2015)

# DocNADE for document representation

* [DocNADE](http://proceedings.mlr.press/v15/larochelle11a/larochelle11a.pdf) and [post how to use it](http://blog.aylien.com/tensorflow-implementation-neural-autoregressive-topic-model-docnade/) to get document representation.

# Deep learning way

Another way to solve similar task

* [Deep API learning, Xiaodong Gu, Hongyu Zhang, Dongmei Zhang, and Sunghun Kim](https://arxiv.org/pdf/1605.08535v1.pdf)
    > We propose DeepAPI, a deep learning based approach to generate API usage sequences for a given natural language query. DeepAPI adapts a neural language model named RNN Encoder-Decoder. It encodes a word sequence (user query) into a fixed-length context vector, and generates an API sequence based on the context vector. We also augment the RNN Encoder-Decoder by considering the importance of individual APIs. We empirically evaluate our approach with more than 7 million annotated code snippets collected from GitHub. The results show that our approach generates largely accurate API sequences and outperforms the related approaches.

    * RNN Encoder-Decoder
        > the API learning problem as a machine translation problem: given a natural language query x = x1 , ..., xn where xt is a key word, we aim to translate it into an API sequence y = y1,...,ym where yt is an API
    * Learns on 
        > a corpus of annotated API sequences, i.e., ⟨API sequence, annotation⟩ pairs,

# Probabilistic way. Interesting Sequence Mining

* [MAST-group.](https://mast-group.github.io/) Probabilistic API Miner. [github](https://github.com/mast-group/api-mining), article [Parameter-Free Probabilistic API Mining across GitHub](https://arxiv.org/pdf/1512.05558.pdf).
    > PAM is a near parameter-free probabilistic algorithm for mining the most interesting API patterns from a list of API call sequences. PAM largely avoids returning redundant and spurious sequences, unlike API mining approaches based on frequent pattern mining.
    
    > the hand-written examples actually have limited coverage of real API usages.

* [MAST-group.](https://mast-group.github.io/) Interesting Sequence Miner [github](https://github.com/mast-group/sequence-mining), article [A Subsequence Interleaving Model for Sequential Pattern Mining](https://arxiv.org/pdf/1602.05012.pdf)
    > ISM is a novel algorithm that mines the subsequences that are most interesting under a probabilistic model of a sequence database. Our model is able to efficiently infer interesting subsequences directly from the database.

### Related works from articles
* MAPO
    * Use code search engines for finding snippets.
    * Clustering by method according to a distance metric, computed as an average of the similarity of method names, class names, and the called API methods themselves.
    * For each cluster, MAPO mines the most frequent API calls using SPAM
* UP-Miner
    * extends MAPO
    * use BIDE algorithm. It returns only the frequent sequences that have no subsequences with the same frequency 
    * clustering distance metric based on the set of all API call sequence n-grams
* SNIFF
    * finds abstract code examples relevant to a natural language query expressing a desired task. 
    * annotates publicly available source code with API documentation and the annotated code is then indexed for searching.

# Can be useful

* [Searching and Skimming: An Exploratory Study. Jamie Starke, Chris Luce, Jonathan Sillito University of Calgary
Calgary, Canada](http://people.ucalgary.ca/~sillito/work/icsm2009.pdf)
    > We conducted a formative study in which programmers were asked to perform corrective tasks to a system they were initially unfamiliar with. Our analysis focused specifically on how programmers decide what to search for, and how they decide which results are relevant to their task. Based on our analysis, we present five observations about our participant’s approach to finding information and some of the challenges they faced. We also discuss the implications these observations have for the design of source code search tools.

# Others articles

* [Code Search via Topic-Enriched Dependence Graph Matching](http://ink.library.smu.edu.sg/cgi/viewcontent.cgi?article=2396&context=sis_research)
    > In this paper, we propose a semantic dependence search engine that integrates both kinds of techniques and can retrieve code snippets based on expressive user queries describing both topics and dependencies. Users can specify their search targets in a free form format describing desired topics (i.e., high-level semantic or functionality of the target code); a specialized graph query language allows users to describe low-level data and control dependencies in code and thus helps to refine the queries described in the free format. Our empirical evaluation on a number of software maintenance tasks shows that our search engine can efficiently locate desired code fragments accurately.
    
    * Use LDA and structural and semantic representations (system dependence graphs (SDGs)) of code.
    * Have special and overcomplicated query language

* [Improving Topic Model Source Code Summarization Paul W. McBurney, Cheng Liu, Collin McMillan, and Tim Weninger](https://www3.nd.edu/~cmc/papers/mcburney_icpcera_2014.pdf)
    > In this paper, we present an emerging source code summarization technique that uses topic modeling to select keywords and topics as summaries for source code. Our approach organizes the topics in source code into a hierarchy, with more general topics near the top of the hierarchy. In this way, we present the software’s highest-level functionality first, before lower-level details. This is an advantage over previous approaches based on topic models, that only present groups of related keywords without a hierarchy. We conducted a preliminary user study that found our approach selects keywords and topics that the participants found to be accurate in a majority of cases.

    * Topic modeling for code summarization task. Task of creating a brief description of a section of source code.
    
    * It uses HDTM algorithm, which analyses graph of documents. 
     > Our technique employs the HDTM algorithm described by Weninger et. al [29] to extract a topic hierarchy for a software system, then we display the hierarchy to programmers in a navigable web interface.

    * We can use this approach to extract names and keywords in themes
    > Code summarization techniques based on topic models are described extensively in software engineering literature [19]. But as a recent study by Panichella et. al points out, these techniques often “have rather low performance when applied on software data” [19].
    
* [API usage pattern recommendation for software development](http://www.sciencedirect.com/science/article/pii/S0164121216301200)
    > Our approach represents the source code as a network of object usages where an object usage is a set of method calls invoked on a single API class. We automatically extract usage patterns by clustering the data based on the co-existence relations between object usages. We conduct an empirical study using a corpus of 11,510 Android applications. The results demonstrate that our approach can effectively mine API usage patterns with high completeness and low redundancy. We observe 18% and 38% improvement on F-measure and response time respectively comparing to usage pattern extraction using frequent-sequence mining.
    
    > Specifically, we show that our approach outperforms the baseline in mining less frequently used API usage patterns. In addition, the ranking quality of our approach is better than Codota which is an online commercial usage pattern recommendation service for Android development.
    
* [Multimodal Code Search by Shaowei Wang](http://ink.library.smu.edu.sg/cgi/viewcontent.cgi?article=1118&context=etd_coll)
    > In this dissertation, we propose a multimodal code search engine, which employs novel techniques that allow developers to effectively find code elements of interest by processing developers’ inputs in various input forms including free-form texts, an SQL-like domain-specific language, code examples, execution traces, and user feedback.  Our evaluations show that our approaches improve over state-of-the-art approaches significantly.
    
    It is a big work. Interesting Literature Review. I think it can be helpful. 

# Libs & Tools

* https://libraries.io/data

    Can be useful to get libraries list and projects, that use specific library.
