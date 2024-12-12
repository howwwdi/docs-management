// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract DocumentManagement {
    struct Document {
        string ipfsHash;
        string fileName;
        uint256 timestamp;
        address owner;
    }

    mapping(uint256 => Document) private documents;

    event DocumentAdded(
        uint256 indexed docId,
        string ipfsHash,
        string fileName,
        address indexed owner,
        uint256 timestamp
    );

    modifier onlyOwner(uint256 docId) {
        require(
            documents[docId].owner == msg.sender,
            "You are not the owner of this document"
        );
        _;
    }

    uint256 private docCounter;

    constructor() {
        docCounter = 0;
    }

    /// @dev func for adding a new document
    /// @param ipfsHash ipfs hash of the document
    /// @param fileName document file name
    /// @return docId uniq document id
    function addDocument(
        string memory ipfsHash,
        string memory fileName
    ) public returns (uint256 docId) {
        docId = docCounter++;
        documents[docId] = Document({
            ipfsHash: ipfsHash,
            fileName: fileName,
            timestamp: block.timestamp,
            owner: msg.sender
        });

        emit DocumentAdded(
            docId,
            ipfsHash,
            fileName,
            msg.sender,
            block.timestamp
        );
    }

    /// @dev func for getting document by id
    /// @param docId document id
    /// @return ipfsHash ipfs hash
    /// @return fileName filename
    /// @return timestamp timestamp
    /// @return owner owner address
    function getDocument(
        uint256 docId
    )
        public
        view
        returns (
            string memory ipfsHash,
            string memory fileName,
            uint256 timestamp,
            address owner
        )
    {
        Document memory doc = documents[docId];
        return (doc.ipfsHash, doc.fileName, doc.timestamp, doc.owner);
    }
 
    /// @dev delete document by id
    /// @param docId uniq document id
    function deleteDocument(uint256 docId) public onlyOwner(docId) {
        delete documents[docId];
    }

    /// @dev get document count
    /// @return count document count
    function getDocumentCount() public view returns (uint256 count) {
        count = docCounter;
    }
}
