# HLF-Library

## Before you begin

We will run the auction smart contract using the Fabric test network. Open a command terminal and navigate to the fabric-samples directory:
```
cd fabric-samples
```

You can then run the following command to deploy the test network. Make sure to set the name as libary.
```
git clone https://github.com/jinwoo0103/hlf-library.git library
```

You should have `library` directory under `fabric-samples` directory now.

If you didn't set environemnt variables for golang, you should set it. You should check golang path rather than just copy and paste below command.
```
export PATH=$PATH:/usr/local/go/bin
```

You can check with following command
```
go version
```

Now, move to library directory.
```
cd library
```

## Start library

You can use simple script file `startLibrary.sh` to start the library.
```
./startLibrary.sh
```

## Library functions
You can perform below 4 actions.

### GetAllAssets
Prints all informations of books.
```
node application-javascript/getAllAssets.js
```

### Purchase
Purchase a book specified using bookID.
```
node application-javascript/purchaseBook.js bookID
```

### Borrow
Borrow a book specified using bookID, together with borrower.
```
node application-javascript/borrowBook.js bookID borrower
```

### Return
Return a book specified using bookID, together with borrower.
```
node application-javascript/returnBook.js bookID borrower
```
