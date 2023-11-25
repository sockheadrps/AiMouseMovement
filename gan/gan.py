import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import Dataset
import json


class MouseMovementDataset(Dataset):
    def __init__(self, file_path):
        self.data = self.read_json_file(file_path)

    def __len__(self):
        return len(self.data)

    def __getitem__(self, idx):
        sample = self.data[idx]
        mouse_array = sample['mousearray']
        x = [point['x'] for point in mouse_array]
        y = [point['y'] for point in mouse_array]
        time = [point['time'] for point in mouse_array]

        return {'x': torch.tensor(x), 'y': torch.tensor(y), 'time': torch.tensor(time)}

    def read_json_file(self, file_path):
        with open(file_path, 'r') as file:
            data = json.load(file)
        return data


# Creating an instance of the dataset
data_file = 'output.json'
dataset = MouseMovementDataset(data_file)


class Generator(nn.Module):
    def __init__(self, input_dim, output_dim):
        super(Generator, self).__init__()
        self.fc = nn.Linear(input_dim, output_dim)

    def forward(self, x):
        return self.fc(x)


class Discriminator(nn.Module):
    def __init__(self, input_dim):
        super(Discriminator, self).__init__()
        self.fc = nn.Linear(input_dim, 1)
        self.sigmoid = nn.Sigmoid()

    def forward(self, x):
        x = self.fc(x)
        return self.sigmoid(x)


# Set parameters
input_dim = 3  # Adjust the input dimensions based on your dataset
output_dim = 3  # Adjust the output dimensions based on your dataset
lr = 0.0002
epochs = 10000
batch_size = 64

# Initialize the models
generator = Generator(input_dim, output_dim)
discriminator = Discriminator(input_dim)

# Define the optimizers
optimizer_G = optim.Adam(generator.parameters(), lr=lr)
optimizer_D = optim.Adam(discriminator.parameters(), lr=lr)

# Assuming 'dataset' is your formatted MouseMovementDataset
dataloader = torch.utils.data.DataLoader(
    dataset, batch_size=batch_size, shuffle=True)

# Training loop
for epoch in range(epochs):
    for i, data in enumerate(dataloader):
        real_data = {'x': data['x'], 'y': data['y'], 'time': data['time']}
        x_tensor = torch.tensor(real_data['x'], dtype=torch.float).view(-1, 1)
        y_tensor = torch.tensor(real_data['y'], dtype=torch.float).view(-1, 1)
        time_tensor = torch.tensor(
            real_data['time'], dtype=torch.float).view(-1, 1)
        real_data = torch.cat((x_tensor, y_tensor, time_tensor), dim=1)

        # Train the discriminator
        optimizer_D.zero_grad()
        criterion = nn.BCELoss()

        real_decision = discriminator(real_data)
        real_error = criterion(real_decision, torch.ones_like(real_decision))

        # Real data
        real_output = discriminator(real_data)
        real_loss = criterion(real_output, torch.ones_like(real_output))

        # Fake data
        fake_data = generator(torch.randn(batch_size, input_dim))
        fake_output = discriminator(fake_data)
        fake_loss = criterion(fake_output, torch.zeros_like(fake_output))

        d_loss = real_loss + fake_loss
        d_loss.backward()
        optimizer_D.step()

        d_loss = real_loss + fake_loss

        # Train the generator
        optimizer_G.zero_grad()
        fake_data = generator(torch.randn(batch_size, input_dim))
        fake_decision = discriminator(fake_data)
        fake_error = criterion(fake_decision, torch.zeros_like(fake_decision))

        g_loss = fake_error  # Define the g_loss here
        g_loss.backward()
        optimizer_G.step()

    # Print the progress
    print(
        f"[Epoch {epoch}/{epochs}] Discriminator Loss: {d_loss.item()} Generator Loss: {g_loss.item()}")


# Use the generator to generate new mouse movement data points
with torch.no_grad():
    # Adjust the input based on your dataset
    fake_data = generator(torch.randn(3))

num_samples = 100  # Define the number of samples you want to generate
generated_data = []
with torch.no_grad():
    for _ in range(num_samples):
        # Adjust the input based on your dataset
        fake_data = generator(torch.randn(1, input_dim))
        generated_data.append(fake_data)

reverted_data = []
for sample in generated_data:
    # Check if the sample has at least three elements along the second dimension
    if sample.shape[1] >= 3:
        x = sample[0][0].item()  # Extract the x coordinate
        y = sample[0][1].item()  # Extract the y coordinate
        time = sample[0][2].item()  # Extract the time value
        reverted_data.append({'x': x, 'y': y, 'time': time})

print(reverted_data)


def write_to_json(data, output_file):
    with open(output_file, 'w') as file:
        json.dump(data, file, indent=2)


# Example usage:
data_structure = [
    {
        "_id": "65568654817d238a6665b2ef",
        "windowheight": 800,
        "windowwidth": 800,
        "mousearray": reverted_data
    },
    # Add more items as needed
]

output_file = 'output.json'
write_to_json(data_structure, output_file)
