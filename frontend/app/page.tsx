"use client";

import React from "react";


export default function Home() {

  const [name, setName] = React.useState("");
  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");

  return (
    <div className="flex flex-col min-h-screen items-center bg-zinc-50 font-sans dark:bg-zinc-950 p-6 sm:p-12 text-zinc-900 dark:text-zinc-50">
      <main className="w-full max-w-4xl bg-white dark:bg-zinc-900 shadow-sm rounded-xl border border-zinc-200 dark:border-zinc-800 p-6 sm:p-8 space-y-8">
        
      
        {/* Header */}
        <div>
          <h1 className="text-xl font-semibold tracking-tight">User Management</h1>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          
          {/* Left Column: Form Wrapper */}
          <div className="md:col-span-1 space-y-4">
            <h2 className="text-sm font-medium uppercase tracking-wider text-zinc-400">
              {/* TODO: Dynamically toggle between "Create User" and "Update User" */}
              Create User
            </h2>
            
            {/* TODO: Bind onSubmit handler */}
            <form className="space-y-4">
              <div>
                <label className="block text-xs font-medium mb-1 text-zinc-600 dark:text-zinc-400">Full Name</label>
                {/* TODO: Add name, value, and onChange props */}
                <input
                  type="text"
                  name="name"
                  value=""
                  placeholder="John Doe"
                  className="w-full px-3 py-2 text-sm bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>

              <div>
                <label className="block text-xs font-medium mb-1 text-zinc-600 dark:text-zinc-400">Email Address</label>
                {/* TODO: Add name, value, and onChange props */}
                <input
                  type="email"
                  placeholder="john@example.com"
                  className="w-full px-3 py-2 text-sm bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>

              {/* TODO: Render conditionally only when creating a new user */}
              <div>
                <label className="block text-xs font-medium mb-1 text-zinc-600 dark:text-zinc-400">Password</label>
                {/* TODO: Add name, value, and onChange props */}
                <input
                  type="password"
                  placeholder="••••••••"
                  className="w-full px-3 py-2 text-sm bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>

              <div className="flex gap-2 pt-2">
                <button
                  type="submit"
                  className="flex-1 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium py-2 px-4 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-zinc-900"
                >
                  {/* TODO: Dynamically toggle between "Add User" and "Save Changes" */}
                  Add User
                </button>
                
                {/* TODO: Render conditionally only when editing an existing user */}
                <button
                  type="button"
                  className="bg-zinc-100 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 text-sm font-medium py-2 px-3 rounded-lg transition-colors"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>

          {/* Right Column: List Wrapper */}
          <div className="md:col-span-2 space-y-4">
            <h2 className="text-sm font-medium uppercase tracking-wider text-zinc-400">
              {/* TODO: Bind dynamic item length count */}
              Active Users (2)
            </h2>

            {/* TODO: Conditional rendering template */}
            {/* IF user array is empty, show this empty state view: */}
            {false ? (
              <div className="text-center py-12 border border-dashed border-zinc-200 dark:border-zinc-800 rounded-xl">
                <p className="text-sm text-zinc-400">No users found. Try adding one on the left!</p>
              </div>
            ) : (
              /* ELSE render the table with data: */
              <div className="overflow-x-auto border border-zinc-200 dark:border-zinc-800 rounded-xl">
                <table className="w-full text-left border-collapse text-sm">
                  <thead>
                    <tr className="bg-zinc-50 dark:bg-zinc-800/50 border-b border-zinc-200 dark:border-zinc-800 text-zinc-500 dark:text-zinc-400 font-medium">
                      <th className="p-4">User</th>
                      <th className="p-4 text-right">Actions</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-zinc-200 dark:divide-zinc-800">
                    {/* TODO: Use .map() to loop through active users here */}
                    {/* Placeholder row start */}
                    <tr className="hover:bg-zinc-50/50 dark:hover:bg-zinc-800/30 transition-colors">
                      <td className="p-4">
                        <div className="font-medium text-zinc-900 dark:text-zinc-100">Alex Rivera</div>
                        <div className="text-xs text-zinc-500 dark:text-zinc-400">alex@example.com</div>
                      </td>
                      <td className="p-4 text-right space-x-2">
                        <button className="text-xs font-medium text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300">
                          Edit
                        </button>
                        <button className="text-xs font-medium text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300">
                          Delete
                        </button>
                      </td>
                    </tr>
                    {/* Placeholder row end */}
                  </tbody>
                </table>
              </div>
            )}
          </div>

        </div>
      </main>
    </div>
  );
}