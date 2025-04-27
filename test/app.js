import fs from "fs";

function main() {
  const filePath = "./debug-test/test.txt";
  const fileContent = fs.readFileSync(filePath, "utf8");

  console.log("File content:", fileContent);

  var data = fileContent.split("");

  var tree = encodeTree(data);

  //console.log("Tree:", JSON.stringify(tree, null, 2));

  var tree = setValuesTree(tree, data);

  // TODO: optimize the tree structure

  var tree = optimizeTree(tree);

  //console.log("Tree with values:", JSON.stringify(tree, null, 2));

  fs.writeFileSync("tree.json", JSON.stringify(tree));

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
          //v: "", //[data[i + 1]],
          //n: {},
        };
      } else {
        //branch[char].v.push(data[i + 1]);
      }

      if (index > 0) {
        if (branch[data[index]].n === undefined) {
          branch[data[index]].n = {};
        }
        expandNode(branch[data[index]].n, index - 1);
      }
    }

    expandNode(tree, i);
  }

  return tree;
}

function setValuesTree(tree, data) {
  function traverse(branch, i, value) {
    if (branch.n === undefined || branch.n[data[i]] === undefined) {
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

function optimizeTree(tree) {
  function traverse(branch, rootBranch) {
    if (branch.n === undefined) {
      rootBranch.v = branch.v;
      rootBranch.n = undefined;
      return;
    }

    if (Object.keys(branch.n).length === 1) {
      traverse(branch.n[Object.keys(branch.n)[0]], rootBranch);
      return;
    } else {
      for (var key in branch.n) {
        traverse(branch.n[key], branch.n[key]);
      }
    }

    var values = {};
    for (var key in branch.n) {
      if (branch.n[key].v !== undefined && branch.n[key].n === undefined) {
        if (values[branch.n[key].v] === undefined) {
          values[branch.n[key].v] = 0;
        }
        values[branch.n[key].v] += 1;
      }
    }

    var max = 0;
    var maxKey = "";
    for (var key in values) {
      if (values[key] > max) {
        max = values[key];
        maxKey = key;
      }
    }

    branch.v = maxKey;
    for (var key in branch.n) {
      if (
        branch.n[key].n === undefined &&
        branch.n[key].v !== undefined &&
        branch.n[key].v === maxKey
      ) {
        branch.n[key] = undefined;
      }
    }
  }

  for (var key in tree) {
    traverse(tree[key], tree[key]);
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
