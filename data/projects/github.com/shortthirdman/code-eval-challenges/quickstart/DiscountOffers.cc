/*
Discount Offers  Share on LinkedIn

Description:

Our marketing department has just negotiated a deal with several local merchants that will allow us to offer exclusive discounts on various products to our top customers every day. The catch is that we can only offer each product to one customer and we may only offer one product to each customer.

Each day we will get the list of products that are eligible for these special discounts. We then have to decide which products to offer to which of our customers. Fortunately, our team of highly skilled statisticians has developed an amazing mathematical model for determining how likely a given customer is to buy an offered product by calculating what we call the "suitability score" (SS). The top-secret algorithm to calculate the SS between a customer and a product is this:

1. If the number of letters in the product's name is even then the SS is the number of vowels (a, e, i, o, u, y) in the customer's name multiplied by 1.5.
2. If the number of letters in the product's name is odd then the SS is the number of consonants in the customer's name.
3. If the number of letters in the product's name shares any common factors (besides 1) with the number of letters in the customer's name then the SS is multiplied by 1.5.

Your task is to implement a program that assigns each customer a product to be offered in a way that maximizes the combined total SS across all of the chosen offers. Note that there may be a different number of products and customers. You may include code from external libraries as long as you cite the source.

Input sample:

Your program should accept as its only argument a path to a file. Each line in this file is one test case. Each test case will be a comma delimited set of customer names followed by a semicolon and then a comma delimited set of product names. Assume the input file is ASCII encoded. For example (NOTE: The example below has 3 test cases):


Jack Abraham,John Evans,Ted Dziuba;iPad 2 - 4-pack,Girl Scouts Thin Mints,Nerf Crossbow
Jeffery Lebowski,Walter Sobchak,Theodore Donald Kerabatsos,Peter Gibbons,Michael Bolton,Samir Nagheenanajar;Half & Half,Colt M1911A1,16lb bowling ball,Red Swingline Stapler,Printer paper,Vibe Magazine Subscriptions - 40 pack
Jareau Wade,Rob Eroh,Mahmoud Abdelkader,Wenyi Cai,Justin Van Winkle,Gabriel Sinkin,Aaron Adelson;Batman No. 1,Football - Official Size,Bass Amplifying Headphones,Elephant food - 1024 lbs,Three Wolf One Moon T-shirt,Dom Perignon 2000 Vintage
Output sample:

For each line of input, print out the maximum total score to two decimal places. For the example input above, the output should look like this:

21.00
83.50
71.25
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <iomanip>
#include <cstdio>
#include <memory.h>


using namespace std;

bool isLetter(char c) {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z');
}

bool isVowel(char c) {
    static string vowels = "aeiouyAEIOUY";
    return (vowels.find(c) != string::npos);
}

int numberOfLetters(string s) {
    int res = 0;
    for (int i = 0; i < s.size(); ++i) {
        if (isLetter(s[i])) ++res;
    }
    return res;
}

int numberOfVowels(string s) {
    int res = 0;
    for (int i = 0; i < s.size(); ++i) {
        if (isVowel(s[i])) ++res;
    }
    return res;
}

int commonFactor(int a, int b) {
    int m = min(a, b);
    for (int i = 2; i <= m; ++i) {
        if (a%i == 0 && b%i == 0)
            return i;
    }
    return 1;
}

double score(string cust, string prod) {    
    double res = 0;
    if (numberOfLetters(prod)%2 == 0) {
        res = (double)numberOfVowels(cust) * 1.5;
    } else {
        res = numberOfLetters(cust) - numberOfVowels(cust);
    }

    if (commonFactor(numberOfLetters(prod), numberOfLetters(cust)) > 1) {
        res *= 1.5;
    }

    return res;
}

void parse(vector<string> &cust, vector<string> &prod, string buf) {
    int pos = buf.find(';');
    string cust_str = buf.substr(0, pos);
    string prod_str = buf.substr(pos+1);

    // generate customer
    pos = cust_str.find(',');
    while (pos != string::npos) {
        cust.push_back(cust_str.substr(0, pos));
        cust_str = cust_str.substr(pos+1);
        pos = cust_str.find(',');
    }
    cust.push_back(cust_str);

    // generate product
    pos = prod_str.find(',');
    while (pos != string::npos) {
        prod.push_back(prod_str.substr(0, pos));
        prod_str = prod_str.substr(pos+1);
        pos = prod_str.find(',');
    }
    prod.push_back(prod_str);
}



// http://files.myopera.com/dd-usaco/cpp/ural1076.cpp?1363818034
// http://blog.csdn.net/yulin11/article/details/4385207
const int MAX = 100;
int n;                // X 的大小
int weight[MAX][MAX];        // X 到 Y 的映射（权重）
int lx[MAX], ly[MAX];        // 标号
bool sx[MAX], sy[MAX];    // 是否被搜索过
int match[MAX];        // Y(i) 与 X(match [i]) 匹配


// 从 X(u) 寻找增广道路，找到则返回 true
bool path (int u);
// 参数 maxsum 为 true ，返回最大权匹配，否则最小权匹配
int bestmatch (bool maxsum = true);
void init (vector<vector<double> > &table)
{
    // 根据实际情况，添加代码以初始化
    n = max(table.size(), table[0].size());
    for (int i = 0; i < n; i ++) {
        for (int j = 0; j < n; j ++) {
            if (i >= table.size()) weight[i][j] = 0;
            else if (j >= table[0].size()) weight[i][j] = 0;
            else weight[i][j] = table[i][j] * 100;
        }
    }
}

bool path (int u)
{
    sx [u] = true;
    for (int v = 0; v < n; v ++)
        if (!sy [v] && lx[u] + ly [v] == weight [u] [v])
            {
            sy [v] = true;
            if (match [v] == -1 || path (match [v]))
                {
                match [v] = u;
                return true;
                }
            }
    return false;
}
int bestmatch (bool maxsum)
{
    int i, j;
    if (!maxsum)
        {
        for (i = 0; i < n; i ++)
            for (j = 0; j < n; j ++)
                weight [i] [j] = -weight [i] [j];
        }
    // 初始化标号
    for (i = 0; i < n; i ++)
        {
        lx [i] = -0x1FFFFFFF;
        ly [i] = 0;
        for (j = 0; j < n; j ++)
            if (lx [i] < weight [i] [j])
                lx [i] = weight [i] [j];
        }
    memset (match, -1, sizeof (match));
    for (int u = 0; u < n; u ++)
        while (1)
            {
            memset (sx, 0, sizeof (sx));
            memset (sy, 0, sizeof (sy));
            if (path (u))
                break;
            // 修改标号
            int dx = 0x7FFFFFFF;
            for (i = 0; i < n; i ++)
                if (sx [i])
                    for (j = 0; j < n; j ++)
                        if(!sy [j])
                            dx = min (lx[i] + ly [j] - weight [i] [j], dx);
            for (i = 0; i < n; i ++)
                {
                if (sx [i])
                    lx [i] -= dx;
                if (sy [i])
                    ly [i] += dx;
                }
            }
    int sum = 0;
    for (i = 0; i < n; i ++)
        sum += weight [match [i]] [i];
    if (!maxsum)
        {
        sum = -sum;
        for (i = 0; i < n; i ++)
            for (j = 0; j < n; j ++)
                weight [i] [j] = -weight [i] [j];         // 如果需要保持 weight [ ] [ ] 原来的值，这里需要将其还原
        }
    return sum;
}



int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string buf;

    while (getline(ifs, buf)) {
        vector<string> cust;
        vector<string> prod;

        parse(cust, prod, buf);
            
        // cal SS value for each possible match
        vector<vector<double> > table(cust.size(), vector<double> (prod.size()));
        for (int i = 0; i < cust.size(); ++i) {
            for (int j = 0; j < prod.size(); ++j) {
                table[i][j] = score(cust[i], prod[j]);
            }
        }

        init(table);
        cout << fixed << setprecision(2) << bestmatch(true)/100.0 << endl;
    }

    return 0;
}