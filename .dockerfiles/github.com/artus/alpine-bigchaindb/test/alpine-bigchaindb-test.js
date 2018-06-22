/**
 * @author Artus Vranken
 */


// Initialize connection
const driver = require('bigchaindb-driver');
const bip39 = require('bip39');
const API_PATH = 'http://localhost:59984/api/v1/';
const conn = new driver.Connection(API_PATH);

// Create identities
const alice = new driver.Ed25519Keypair(bip39.mnemonicToSeed("alice").slice(0, 32));
const bob = new driver.Ed25519Keypair(bip39.mnemonicToSeed("bob").slice(0, 32));

const assert = require('assert');

describe('artusvranken/alpine-bigchaindb docker image', function () {

    describe('Create transactions (adding assets)', function () {

        it('should add assets when done correctly', function (done) {
            this.timeout(10000);

            // Create a new asset.
            const assetData = {
                'testAsset': 'This is a test-asset.'
            };

            // Create metadata.
            const metaData = {
                'date': new Date()
            }

            // Create transaction.
            const unsignedCreateTransaction = driver.Transaction.makeCreateTransaction(
                assetData,
                metaData,
                [driver.Transaction.makeOutput(driver.Transaction.makeEd25519Condition(alice.publicKey))],
                alice.publicKey
            );

            // Sign transaction.
            const signedTransaction = driver.Transaction.signTransaction(unsignedCreateTransaction, alice.privateKey);

            // Post transaction.
            conn.postTransaction(signedTransaction).then(response => {

                // Poll for status.
                return conn.pollStatusAndFetchTransaction(response.id);
            }).then(whatever => {

                // Check if the asset was posted.
                return conn.listOutputs(alice.publicKey);
            }).then(response => {

                // We do so by checking if there is exactly one transaction.
                if (response.length == 1) done();
                else (done(new Error("Length was not 1.")));
            }).catch(error => {
                done(new Error(error));
            })
        });

        it('should fail when not done correctly', function (done) {
            this.timeout(10000);

            // Create a new asset.
            const assetData = "incorrect";

            // Create metadata.
            const metaData = "incorrect";

            // Create transaction.
            const unsignedCreateTransaction = driver.Transaction.makeCreateTransaction(
                assetData,
                metaData,
                [driver.Transaction.makeOutput(driver.Transaction.makeEd25519Condition(alice.publicKey))],
                alice.publicKey
            );

            // Sign transaction
            const signedTransaction = driver.Transaction.signTransaction(unsignedCreateTransaction, alice.privateKey);

            // Post transaction.
            conn.postTransaction(signedTransaction).then(response => {

                // Poll for status.
                return conn.pollStatusAndFetchTransaction(response.id);
            }).then(whatever => {

                // Check if the asset was posted.
                return conn.listOutputs(alice.publicKey);
            }).then(response => {

                // We do so by checking if there is exactly one transaction.
                if (response.length == 2) done(new Error("incorrect transaction was added."));
                else (done());
            }).catch(error => {
                done();
            })
        });
    });

    describe('Transfer transactions', function () {

        it('should transfer assets when done correctly', function (done) {
            this.timeout(10000);

            // Create a new asset.
            const assetData = {
                'testAsset': 'This is a test-asset.'
            };

            // Create metadata.
            const metaData = {
                'date': new Date()
            }

            // Create transaction.
            const unsignedCreateTransaction = driver.Transaction.makeCreateTransaction(
                assetData,
                metaData,
                [driver.Transaction.makeOutput(driver.Transaction.makeEd25519Condition(alice.publicKey))],
                alice.publicKey
            );

            // Sign transaction.
            const signedTransaction = driver.Transaction.signTransaction(unsignedCreateTransaction, alice.privateKey);

            // Post transaction.
            conn.postTransaction(signedTransaction).then(response => {

                // Poll for status.
                return conn.pollStatusAndFetchTransaction(response.id);
            }).then(whatever => {

                // Retrieve the latest transaction performed by alice.
                return conn.listOutputs(alice.publicKey, false);
            }).then(response => {


                const latestTransactionOutput = response[response.length - 1];

                return conn.getTransaction(latestTransactionOutput.transaction_id);
            }).then(latestTransaction => {


                // perform transfer transaction on latest asset.
                const transferTransaction = driver.Transaction.makeTransferTransaction(
                    [{ tx: latestTransaction, output_index: 0 }],
                    [driver.Transaction.makeOutput(driver.Transaction.makeEd25519Condition(bob.publicKey))],
                    { 'action': 'transfer to bob' }
                );


                // Sign the transfer transaction.
                const signedTransferTransaction = driver.Transaction.signTransaction(transferTransaction, alice.privateKey);

                return conn.postTransaction(signedTransferTransaction);
            }).then(postedTransaction => {

                return conn.pollStatusAndFetchTransaction(postedTransaction.id);
            }).then(response => {

                done();
            }).catch(error => {
                done(new Error(error));
            });

        });

        it('should fail when done incorrectly', function () {
            this.timeout(10000);

            // Create a new asset.
            const assetData = {
                'testAsset': 'This is a test-asset.'
            };

            // Create metadata.
            const metaData = {
                'date': new Date()
            }

            // Create transaction.
            const unsignedCreateTransaction = driver.Transaction.makeCreateTransaction(
                assetData,
                metaData,
                [driver.Transaction.makeOutput(driver.Transaction.makeEd25519Condition(alice.publicKey))],
                alice.publicKey
            );

            // Sign transaction.
            const signedTransaction = driver.Transaction.signTransaction(unsignedCreateTransaction, alice.privateKey);

            // Post transaction.
            conn.postTransaction(signedTransaction).then(response => {

                // Poll for status.
                return conn.pollStatusAndFetchTransaction(response.id);
            }).then(response => {

                // Retrieve the latest transaction performed by alice.
                return conn.listOutputs(alice.publicKey, false);
            }).then(response => {

                const latestTransaction = response[response.length - 1];

                // perform transfer transaction on latest asset.
                const transferTransaction = driver.Transaction.makeTransferTransaction(
                    [{ tx: latestTransaction, output_index: 0 }],
                    [driver.Transaction.makeOutput(driver.Transaction.makeEd25519Condition(bob.publicKey))],
                    { 'action': 'transfer to bob' }
                );

                // Sign the transfer transaction with an incorrect key.
                const signedTransferTransaction = driver.Transaction.signTransaction(transferTransaction, bob.privateKey);

                return conn.postTransaction(signedTransferTransaction);
            }).then(postedTransaction => {

                return conn.pollStatusAndFetchTransaction(postedTransaction.id);
            }).then(response => {

                if (response['outputs'][0]['public_keys'][0] == bob.publicKey)
                    done(new Error("Incorrect transaction was posted."));
                else throw new Error("Bob is not the new owner.");
            }).catch(error => {
                done();
            });
        });

    });
});