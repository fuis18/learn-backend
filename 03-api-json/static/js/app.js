const API = "/api/products";

const form = document.getElementById("product-form");
const inputId = document.getElementById("id");
const inputName = document.getElementById("name");
const inputPrice = document.getElementById("price");
const inputStock = document.getElementById("stock");
const submitButton = form.querySelector('button[type="submit"]');
const message = document.getElementById("message");

function cell(value) {
  const td = document.createElement("td");
  td.textContent = value;
  return td;
}

async function listProducts() {
  const res = await fetch(API);
  const products = await res.json();

  const tbody = document.getElementById("products-body");
  tbody.replaceChildren();

  for (const product of products) {
    const tr = document.createElement("tr");
    tr.append(
      cell(product.id),
      cell(product.name),
      cell(product.price),
      cell(product.stock),
    );

    const actions = document.createElement("td");

    const edit = document.createElement("button");
    edit.type = "button";
    edit.textContent = "Edit";
    edit.addEventListener("click", () => fillForm(product));

    const remove = document.createElement("button");
    remove.type = "button";
    remove.textContent = "Delete";
    remove.addEventListener("click", () => deleteProduct(product.id));

    actions.append(edit, remove);
    tr.append(actions);

    tbody.append(tr);
  }
}

listProducts();

form.addEventListener("submit", async (event) => {
  event.preventDefault();

  const id = inputId.value;
  const body = JSON.stringify({
    name: inputName.value,
    price: Number(inputPrice.value),
    stock: Number(inputStock.value),
  });

  const res = await fetch(id ? `${API}/${id}` : API, {
    method: id ? "PUT" : "POST",
    headers: { "Content-Type": "application/json" },
    body,
  });

  if (!res.ok) {
    await showError(res);
    return;
  }

  form.reset();
  await listProducts();
});

function fillForm(product) {
  inputId.value = product.id;
  inputName.value = product.name;
  inputPrice.value = product.price;
  inputStock.value = product.stock;
  submitButton.textContent = "Update";
  message.textContent = "";
  inputName.focus();
}

form.addEventListener("reset", () => {
  submitButton.textContent = "Save";
  message.textContent = "";
});

async function deleteProduct(id) {
  const res = await fetch(`${API}/${id}`, { method: "DELETE" });

  if (!res.ok) {
    await showError(res);
    return;
  }

  message.textContent = "";
  await listProducts();
}

// Read {"error": "..."} from failed responses
async function showError(res) {
  let text = `Error ${res.status}`;
  try {
    const body = await res.json();
    if (body.error) {
      text = body.error;
    }
  } catch {
    // body was not JSON
  }
  message.textContent = text;
}
