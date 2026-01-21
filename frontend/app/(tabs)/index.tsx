import React, { useState, useEffect } from 'react';
import {
  StyleSheet,
  Text,
  View,
  TextInput,
  TouchableOpacity,
  FlatList,
  ActivityIndicator,
  Alert,
  SafeAreaView,
} from 'react-native';
import axios from 'axios';

// IMPORTANT: Change this IP address to match your setup
// For Android Emulator: use 10.0.2.2
// For iOS Simulator: use localhost
// For physical device: use your computer's local IP (e.g., 192.168.1.100)
const API_URL = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080/api/tasks';

// TypeScript interface for Task
interface Task {
  id: string;
  title: string;
  description: string;
  completed: boolean;
  created_at: string;
}

export default function TaskManagerScreen() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [loading, setLoading] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);

  // Fetch tasks when the screen loads
  useEffect(() => {
    fetchTasks();
  }, []);

  // GET: Fetch all tasks from the backend
  const fetchTasks = async () => {
    setLoading(true);
    try {
      const response = await axios.get<Task[]>(API_URL);
      setTasks(response.data);
      console.log('Tasks fetched:', response.data);
    } catch (error) {
      console.error('Error fetching tasks:', error);
      Alert.alert('Error', 'Failed to fetch tasks. Make sure your backend is running.');
    } finally {
      setLoading(false);
    }
  };

  // POST: Create a new task
  const createTask = async () => {
    if (!title.trim()) {
      Alert.alert('Error', 'Please enter a task title');
      return;
    }

    try {
      const newTask = {
        title: title,
        description: description,
        completed: false,
      };

      const response = await axios.post<Task>(API_URL, newTask);
      setTasks([...tasks, response.data]);
      setTitle('');
      setDescription('');
      Alert.alert('Success', 'Task created!');
    } catch (error) {
      console.error('Error creating task:', error);
      Alert.alert('Error', 'Failed to create task');
    }
  };

  // PUT: Update an existing task
  const updateTask = async () => {
    if (!editingTask) return;

    try {
      const updatedTask = {
        title: title,
        description: description,
        completed: editingTask.completed,
      };

      const response = await axios.put<Task>(
        `${API_URL}/${editingTask.id}`,
        updatedTask
      );

      setTasks(tasks.map(task =>
        task.id === editingTask.id ? response.data : task
      ));

      setTitle('');
      setDescription('');
      setEditingTask(null);
      Alert.alert('Success', 'Task updated!');
    } catch (error) {
      console.error('Error updating task:', error);
      Alert.alert('Error', 'Failed to update task');
    }
  };

  // DELETE: Delete a task
  const deleteTask = async (taskId: string) => {
    try {
      await axios.delete(`${API_URL}/${taskId}`);
      setTasks(tasks.filter(task => task.id !== taskId));
      Alert.alert('Success', 'Task deleted!');
    } catch (error) {
      console.error('Error deleting task:', error);
      Alert.alert('Error', 'Failed to delete task');
    }
  };

  // Toggle task completion status
  const toggleComplete = async (task: Task) => {
    try {
      const updatedTask = {
        ...task,
        completed: !task.completed,
      };

      const response = await axios.put<Task>(
        `${API_URL}/${task.id}`,
        updatedTask
      );

      setTasks(tasks.map(t =>
        t.id === task.id ? response.data : t
      ));
    } catch (error) {
      console.error('Error toggling task:', error);
      Alert.alert('Error', 'Failed to update task');
    }
  };

  // Start editing a task
  const startEditing = (task: Task) => {
    setEditingTask(task);
    setTitle(task.title);
    setDescription(task.description);
  };

  // Cancel editing
  const cancelEditing = () => {
    setEditingTask(null);
    setTitle('');
    setDescription('');
  };

  // Render a single task item
  const renderTask = ({ item }: { item: Task }) => (
    <View style={styles.taskItem}>
      <TouchableOpacity
        style={styles.taskContent}
        onPress={() => toggleComplete(item)}
      >
        <View style={styles.checkbox}>
          {item.completed && <View style={styles.checkboxChecked} />}
        </View>
        <View style={styles.taskText}>
          <Text style={[
            styles.taskTitle,
            item.completed && styles.completedText
          ]}>
            {item.title}
          </Text>
          {item.description ? (
            <Text style={styles.taskDescription}>{item.description}</Text>
          ) : null}
        </View>
      </TouchableOpacity>

      <View style={styles.taskActions}>
        <TouchableOpacity
          style={styles.editButton}
          onPress={() => startEditing(item)}
        >
          <Text style={styles.editButtonText}>Edit</Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={styles.deleteButton}
          onPress={() => deleteTask(item.id)}
        >
          <Text style={styles.deleteButtonText}>Delete</Text>
        </TouchableOpacity>
      </View>
    </View>
  );

  return (
    <SafeAreaView style={styles.container}>
      <Text style={styles.header}>Task Manager</Text>

      {/* Input Form */}
      <View style={styles.inputContainer}>
        <TextInput
          style={styles.input}
          placeholder="Task Title"
          value={title}
          onChangeText={setTitle}
        />
        <TextInput
          style={styles.input}
          placeholder="Description (optional)"
          value={description}
          onChangeText={setDescription}
        />

        <View style={styles.buttonRow}>
          {editingTask ? (
            <>
              <TouchableOpacity
                style={[styles.button, styles.updateButton]}
                onPress={updateTask}
              >
                <Text style={styles.buttonText}>Update Task</Text>
              </TouchableOpacity>
              <TouchableOpacity
                style={[styles.button, styles.cancelButton]}
                onPress={cancelEditing}
              >
                <Text style={styles.buttonText}>Cancel</Text>
              </TouchableOpacity>
            </>
          ) : (
            <TouchableOpacity
              style={styles.button}
              onPress={createTask}
            >
              <Text style={styles.buttonText}>Add Task</Text>
            </TouchableOpacity>
          )}
        </View>
      </View>

      {/* Task List */}
      {loading ? (
        <ActivityIndicator size="large" color="#007AFF" />
      ) : (
        <FlatList
          data={tasks}
          renderItem={renderTask}
          keyExtractor={item => item.id}
          contentContainerStyle={styles.listContainer}
          ListEmptyComponent={
            <Text style={styles.emptyText}>No tasks yet. Create one above!</Text>
          }
        />
      )}
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  header: {
    fontSize: 28,
    fontWeight: 'bold',
    textAlign: 'center',
    marginTop: 20,
    marginBottom: 20,
    color: '#333',
  },
  inputContainer: {
    backgroundColor: 'white',
    padding: 15,
    marginHorizontal: 15,
    marginBottom: 20,
    borderRadius: 10,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  input: {
    borderWidth: 1,
    borderColor: '#ddd',
    padding: 12,
    marginBottom: 10,
    borderRadius: 8,
    fontSize: 16,
  },
  buttonRow: {
    flexDirection: 'row',
    gap: 10,
  },
  button: {
    backgroundColor: '#007AFF',
    padding: 15,
    borderRadius: 8,
    flex: 1,
  },
  updateButton: {
    backgroundColor: '#34C759',
  },
  cancelButton: {
    backgroundColor: '#FF9500',
  },
  buttonText: {
    color: 'white',
    textAlign: 'center',
    fontSize: 16,
    fontWeight: '600',
  },
  listContainer: {
    paddingHorizontal: 15,
    paddingBottom: 20,
  },
  taskItem: {
    backgroundColor: 'white',
    padding: 15,
    marginBottom: 10,
    borderRadius: 10,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 3,
    elevation: 2,
  },
  taskContent: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 10,
  },
  checkbox: {
    width: 24,
    height: 24,
    borderRadius: 12,
    borderWidth: 2,
    borderColor: '#007AFF',
    marginRight: 12,
    justifyContent: 'center',
    alignItems: 'center',
  },
  checkboxChecked: {
    width: 14,
    height: 14,
    borderRadius: 7,
    backgroundColor: '#007AFF',
  },
  taskText: {
    flex: 1,
  },
  taskTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
  },
  taskDescription: {
    fontSize: 14,
    color: '#666',
    marginTop: 4,
  },
  completedText: {
    textDecorationLine: 'line-through',
    color: '#999',
  },
  taskActions: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    gap: 10,
  },
  editButton: {
    backgroundColor: '#007AFF',
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 6,
  },
  editButtonText: {
    color: 'white',
    fontSize: 14,
    fontWeight: '600',
  },
  deleteButton: {
    backgroundColor: '#FF3B30',
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 6,
  },
  deleteButtonText: {
    color: 'white',
    fontSize: 14,
    fontWeight: '600',
  },
  emptyText: {
    textAlign: 'center',
    color: '#999',
    marginTop: 40,
    fontSize: 16,
  },
});