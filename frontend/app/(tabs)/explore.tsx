import React, { useState } from 'react';
import { SafeAreaView, StyleSheet, TextInput, Text, View } from 'react-native';

const InputOutputExample = () => {
  // 1. Create a state variable to hold the input text
  const [text, setText] = useState('');

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.contentContainer}>

        {/* Label */}
        <Text style={styles.label}>Type something below:</Text>

        {/* 2. The Input Field */}
        <TextInput
          style={styles.input}
          onChangeText={(newText) => setText(newText)} // Updates state when typing
          value={text} // Binds the input value to the state
          placeholder="Enter text here..."
          keyboardType="default"
        />

        {/* 3. The Display Area */}
        <View style={styles.resultBox}>
          <Text style={styles.resultLabel}>Live Output:</Text>
          {/* Display the text state here */}
          <Text style={styles.resultText}>
            {text ? text : "Waiting for input..."}
          </Text>
        </View>

      </View>
    </SafeAreaView>
  );
};

// Styles to make it look clean
const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
    justifyContent: 'center',
  },
  contentContainer: {
    padding: 20,
  },
  label: {
    fontSize: 18,
    fontWeight: '600',
    marginBottom: 10,
    color: '#333',
  },
  input: {
    height: 50,
    borderColor: '#ccc',
    borderWidth: 1,
    borderRadius: 8,
    paddingHorizontal: 15,
    backgroundColor: '#fff',
    fontSize: 16,
    marginBottom: 30,
  },
  resultBox: {
    backgroundColor: '#e3f2fd', // Light blue background
    padding: 20,
    borderRadius: 10,
    alignItems: 'center',
  },
  resultLabel: {
    fontSize: 14,
    color: '#1976d2',
    fontWeight: 'bold',
    marginBottom: 5,
  },
  resultText: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#1565c0',
    textAlign: 'center',
  },
});

export default InputOutputExample;