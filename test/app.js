import fs from "fs";

function main() {
  const filePath = "test.txt";
  const fileContent = fs.readFileSync(filePath, "utf8");

  console.log("File content:", fileContent);

  var tree = encodeTree(fileContent.split(""));

  //console.log("Tree:", JSON.stringify(tree, null, 2));

  var tree = setValuesTree(tree, fileContent.split(""));

  console.log("Tree with values:", JSON.stringify(tree, null, 2));

  // TODO: optimize the tree structure

  var tree1 = {
    l: {
      v: "",
      n: {
        l: {
          v: "o",
          n: {},
        },
        e: {
          v: "l",
          n: {},
        },
      },
    },
    e: {
      v: "l",
      n: {},
    },
    H: {
      v: "e",
      n: {},
    },
  };

  var text = decodeTree(
    tree,
    fileContent.split("")[0],
    fileContent.split("").length
  );

  console.log("Decoded text:", text);
}

function encodeTree(data) {
  var tree = {};

  for (var i = data.length - 2; i >= 0; i--) {
    function expandNode(branch, index) {
      if (branch[data[index]] === undefined) {
        branch[data[index]] = {
          v: "", //[data[i + 1]],
          n: {},
        };
      } else {
        //branch[char].v.push(data[i + 1]);
      }

      if (index > 0) {
        expandNode(branch[data[index]].n, index - 1);
      }
    }

    expandNode(tree, i);
  }

  return tree;
}

function setValuesTree(tree, data) {
  function traverse(branch, i, value) {
    if (branch.n[data[i]] === undefined) {
      branch.v = value;
    } else {
      traverse(branch.n[data[i]], i - 1, value);
    }
  }

  for (var i = data.length - 1; i > 0; i--) {
    traverse(tree[data[i - 1]], i - 2, data[i]);
  }

  return tree;
}

function decodeTree(tree, char, length) {
  var result = [char];

  function traverse(branch, i) {
    if (i < 0 || branch.n === undefined || branch.n[result[i]] === undefined) {
      result.push(branch.v[0]);
    } else {
      traverse(branch.n[result[i]], i - 1);
    }
  }

  for (var i = 0; i < length - 1; i++) {
    traverse(tree[result[i]], i - 1);
  }

  return result.join("");
}

main();
