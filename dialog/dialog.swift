#!/usr/bin/swift

//
//  ReusableUserInputForm.swift
//
//  Created by andy4222 in 2025
//
//  MIT License
//
//  Copyright © 2025 andy4222
//
//  Permission is hereby granted, free of charge, to any person obtaining a copy
//  of this software and associated documentation files (the "Software"), to deal
//  in the Software without restriction, including without limitation the rights
//  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//  copies of the Software, and to permit persons to whom the Software is
//  furnished to do so, subject to the following conditions:
//
//  The above copyright notice and this permission notice shall be included in
//  all copies or substantial portions of the Software.
//
//  ***Attribution is required.*** If you use or adapt this code, you must clearly
//  credit the original author: andy4222 https://www.alfredforum.com/profile/18304-andy4222/
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//  SOFTWARE.
//

import AppKit
import SwiftUI

// MARK: - Command-line Config Loading
func loadConfig(from pathOrJSON: String) -> Config {
    let data: Data

    // Try to parse as JSON string first
    if let jsonData = pathOrJSON.data(using: .utf8),
       let _ = try? JSONDecoder().decode(Config.self, from: jsonData) {
        data = jsonData
    } else {
        // Fallback to file path
        let url = URL(fileURLWithPath: pathOrJSON)
        do {
            data = try Data(contentsOf: url)
        } catch {
            fatalError("Error loading config from \(pathOrJSON): \(error)")
        }
    }

    do {
        let config = try JSONDecoder().decode(Config.self, from: data)
        return config
    } catch {
        fatalError("Error decoding config: \(error)")
    }
}

// MARK: - Config Models
struct Field: Identifiable, Codable, Comparable {
    enum FieldType: String, Codable {
        case text, checkbox, dropdown, texteditor, segmented, filepicker
    }
    let id = UUID()
    let type: FieldType
    let label: String
    let bindingKey: String
    let options: [String]?
    let defaultValue: String?
    let filePickerType: String? // "file" or "folder"
    let copy: Bool? // show copy button for text/texteditor
    let note: String? // note text shown below field in red
    let order: Int // display order

    private enum CodingKeys: String, CodingKey {
        case type, label, bindingKey, options, defaultValue, filePickerType, copy, note, order
    }

    // Implement Comparable for sorting
    static func < (lhs: Field, rhs: Field) -> Bool {
        return lhs.order < rhs.order
    }
}

struct Config: Codable {
    let windowTitle: String?
    let windowWidth: CGFloat?
    let windowHeight: CGFloat?
    let fields: [String: Field]
    let okLabel: String?
    let cancelLabel: String?
    let alwaysOnTop: Bool?

    // Computed property to get fields as sorted array
    var sortedFields: [Field] {
        return Array(fields.values).sorted()
    }
}

// MARK: - NSTextViewRepresentable for Resizable TextEditor (CustomTextView)
final class CustomTextView: NSTextView {
    override func keyDown(with event: NSEvent) {
        if let characters = event.charactersIgnoringModifiers, characters == "\t" {
            if self.string.isEmpty {
                self.window?.selectNextKeyView(self)
                return
            }
        }
        super.keyDown(with: event)
    }
}

struct NSTextViewRepresentable: NSViewRepresentable {
    @Binding var text: String

    func makeCoordinator() -> Coordinator { Coordinator(self) }

    func makeNSView(context: Context) -> NSScrollView {
        let textView = CustomTextView()
        textView.isRichText = false
        textView.isEditable = true
        textView.delegate = context.coordinator
        textView.font = NSFont.systemFont(ofSize: NSFont.systemFontSize)
        textView.backgroundColor = NSColor.textBackgroundColor
        textView.isHorizontallyResizable = true
        textView.isVerticallyResizable = true
        textView.autoresizingMask = [.width]

        if let container = textView.textContainer {
            container.widthTracksTextView = true
            // Explicitly use CGFloat.greatestFiniteMagnitude to avoid ambiguity
            container.containerSize = NSSize(width: 0, height: CGFloat.greatestFiniteMagnitude)
        }

        let scrollView = NSScrollView()
        scrollView.documentView = textView
        scrollView.hasVerticalScroller = true
        scrollView.autohidesScrollers = true
        return scrollView
    }

    func updateNSView(_ nsView: NSScrollView, context: Context) {
        guard let textView = nsView.documentView as? NSTextView else { return }
        if textView.string != text {
            textView.string = text
        }
    }

    final class Coordinator: NSObject, NSTextViewDelegate {
        var parent: NSTextViewRepresentable
        init(_ parent: NSTextViewRepresentable) { self.parent = parent }
        func textDidChange(_ notification: Notification) {
            if let textView = notification.object as? NSTextView {
                parent.text = textView.string
            }
        }
    }
}

// MARK: - Fixed Size TextEditor
struct FixedTextEditor: View {
    @Binding var text: String
    let height: CGFloat

    var body: some View {
        NSTextViewRepresentable(text: $text)
            .frame(height: height)
            .border(Color.gray.opacity(0.2), width: 1)
    }
}

// MARK: - FlowRadioGroup (Horizontal, Wrapping Radio Buttons)
struct FlowRadioGroup: View {
    let options: [String]
    @Binding var selection: String

    var body: some View {
        let columns = [GridItem(.adaptive(minimum: 60), spacing: 8)]
        LazyVGrid(columns: columns, alignment: .leading, spacing: 8) {
            ForEach(options, id: \.self) { option in
                HStack(spacing: 8) {
                    Image(systemName: selection == option ? "largecircle.fill.circle" : "circle")
                    Text(option)
                }
                .onTapGesture {
                    selection = option
                }
            }
        }
        .focusable()
        .onMoveCommand { direction in
            guard let currentIndex = options.firstIndex(of: selection) else { return }
            switch direction {
            case .left, .up:
                let newIndex = (currentIndex - 1 + options.count) % options.count
                selection = options[newIndex]
            case .right, .down:
                let newIndex = (currentIndex + 1) % options.count
                selection = options[newIndex]
            default:
                break
            }
        }
    }
}

// MARK: - FilePickerField (File/Folder Picker)
struct FilePickerField: View {
    @Binding var path: String
    let pickerType: String // "file" or "folder"

    var body: some View {
        HStack(spacing: 8) {
            TextField("", text: $path)
                .textFieldStyle(.roundedBorder)

            Button("浏览...") {
                let panel = NSOpenPanel()
                panel.canChooseFiles = (pickerType == "file")
                panel.canChooseDirectories = (pickerType == "folder")
                panel.allowsMultipleSelection = false

                if panel.runModal() == .OK, let url = panel.url {
                    path = url.path
                }
            }
            .frame(width: 80)
        }
    }
}

// MARK: - DynamicDialogView
struct DynamicDialogView: View {
    let fields: [Field]
    let okLabel: String
    let cancelLabel: String
    @State private var values: [String: Any]

    // Dynamically compute the widest label among all fields.
    private var maxLabelWidth: CGFloat {
        let labels = fields.map { $0.label + ":" }
        let font = NSFont.systemFont(ofSize: NSFont.systemFontSize)
        // Measure the width of each label string
        let widest =
            labels.map { (label: String) -> CGFloat in
                let size = (label as NSString).size(withAttributes: [.font: font])
                return size.width
            }.max() ?? 80
        // Add a little extra padding
        return widest + 8
    }

    init(fields: [Field], okLabel: String, cancelLabel: String) {
        self.fields = fields
        self.okLabel = okLabel
        self.cancelLabel = cancelLabel

        var initialValues: [String: Any] = [:]
        for field in fields {
            if let defaultValue = field.defaultValue {
                switch field.type {
                case .checkbox:
                    initialValues[field.bindingKey] = (defaultValue.lowercased() == "true")
                default:
                    initialValues[field.bindingKey] = defaultValue
                }
            }
        }
        _values = State(initialValue: initialValues)
    }

    var body: some View {
        VStack(spacing: 12) {
            ForEach(fields) { field in
                switch field.type {
                case .text:
                    VStack(alignment: .leading, spacing: 4) {
                        HStack(alignment: .center, spacing: 10) {
                            Text(field.label + ":")
                                .lineLimit(1)
                                .frame(width: maxLabelWidth, alignment: .leading)
                            TextField(
                                "",
                                text: Binding(
                                    get: { self.values[field.bindingKey] as? String ?? "" },
                                    set: { self.values[field.bindingKey] = $0 }
                                )
                            )
                            .textFieldStyle(.roundedBorder)

                            if field.copy == true {
                                Button(action: {
                                    let pasteboard = NSPasteboard.general
                                    pasteboard.clearContents()
                                    pasteboard.setString(
                                        self.values[field.bindingKey] as? String ?? "",
                                        forType: .string
                                    )
                                }) {
                                    Label("", systemImage: "doc.on.doc")
                                }
                                .labelStyle(.iconOnly)
                                .frame(width: 30)
                            } else {
                                Spacer().frame(width: 0)
                            }
                        }
                        if let note = field.note, !note.isEmpty {
                            Text(note)
                                .font(.caption)
                                .foregroundColor(.red)
                                .padding(.leading, maxLabelWidth + 10)
                        }
                    }
                case .checkbox:
                    VStack(alignment: .leading, spacing: 4) {
                        HStack(alignment: .center, spacing: 10) {
                            Text(field.label + ":")
                                .lineLimit(1)
                                .frame(width: maxLabelWidth, alignment: .leading)
                            Toggle(
                                "",
                                isOn: Binding(
                                    get: { self.values[field.bindingKey] as? Bool ?? false },
                                    set: { self.values[field.bindingKey] = $0 }
                                )
                            )
                            .labelsHidden()
                            .toggleStyle(CheckboxToggleStyle())
                            Spacer()
                        }
                        if let note = field.note, !note.isEmpty {
                            Text(note)
                                .font(.caption)
                                .foregroundColor(.red)
                                .padding(.leading, maxLabelWidth + 10)
                        }
                    }
                case .dropdown:
                    if let options = field.options, !options.isEmpty {
                        VStack(alignment: .leading, spacing: 4) {
                            HStack(alignment: .center, spacing: 10) {
                                Text(field.label + ":")
                                    .lineLimit(1)
                                    .frame(width: maxLabelWidth, alignment: .leading)
                                Picker(
                                    "",
                                    selection: Binding(
                                        get: {
                                            self.values[field.bindingKey] as? String ?? options.first!
                                        },
                                        set: { self.values[field.bindingKey] = $0 }
                                    )
                                ) {
                                    ForEach(options, id: \.self) { option in
                                        Text(option).tag(option)
                                    }
                                }
                                .pickerStyle(MenuPickerStyle())
                                Spacer()
                            }
                            if let note = field.note, !note.isEmpty {
                                Text(note)
                                    .font(.caption)
                                    .foregroundColor(.red)
                                    .padding(.leading, maxLabelWidth + 10)
                            }
                        }
                    }
                case .texteditor:
                    VStack(alignment: .leading, spacing: 4) {
                        HStack(alignment: .top, spacing: 10) {
                            Text(field.label + ":")
                                .frame(width: maxLabelWidth, alignment: .leading)
                                .padding(.top, 8)
                            VStack(spacing: 0) {
                                FixedTextEditor(
                                    text: Binding(
                                        get: { self.values[field.bindingKey] as? String ?? "" },
                                        set: { self.values[field.bindingKey] = $0 }
                                    ),
                                    height: 100
                                )
                            }

                            if field.copy == true {
                                Button(action: {
                                    let pasteboard = NSPasteboard.general
                                    pasteboard.clearContents()
                                    pasteboard.setString(
                                        self.values[field.bindingKey] as? String ?? "",
                                        forType: .string
                                    )
                                }) {
                                    Label("", systemImage: "doc.on.doc")
                                }
                                .labelStyle(.iconOnly)
                                .frame(width: 30)
                                .padding(.top, 8)
                            } else {
                                Spacer().frame(width: 0)
                            }
                        }
                        if let note = field.note, !note.isEmpty {
                            Text(note)
                                .font(.caption)
                                .foregroundColor(.red)
                                .padding(.leading, maxLabelWidth + 10)
                        }
                    }
                case .segmented:
                    if let options = field.options, !options.isEmpty {
                        VStack(alignment: .leading, spacing: 4) {
                            HStack(alignment: .center, spacing: 10) {
                                Text(field.label + ":")
                                    .lineLimit(1)
                                    .frame(width: maxLabelWidth, alignment: .leading)
                                FlowRadioGroup(
                                    options: options,
                                    selection: Binding(
                                        get: {
                                            self.values[field.bindingKey] as? String ?? options.first!
                                        },
                                        set: { self.values[field.bindingKey] = $0 }
                                    )
                                )
                                Spacer()
                            }
                            if let note = field.note, !note.isEmpty {
                                Text(note)
                                    .font(.caption)
                                    .foregroundColor(.red)
                                    .padding(.leading, maxLabelWidth + 10)
                            }
                        }
                    }
                case .filepicker:
                    VStack(alignment: .leading, spacing: 4) {
                        HStack(alignment: .center, spacing: 10) {
                            Text(field.label + ":")
                                .lineLimit(1)
                                .frame(width: maxLabelWidth, alignment: .leading)
                            FilePickerField(
                                path: Binding(
                                    get: { self.values[field.bindingKey] as? String ?? "" },
                                    set: { self.values[field.bindingKey] = $0 }
                                ),
                                pickerType: field.filePickerType ?? "file"
                            )
                        }
                        if let note = field.note, !note.isEmpty {
                            Text(note)
                                .font(.caption)
                                .foregroundColor(.red)
                                .padding(.leading, maxLabelWidth + 10)
                        }
                    }
                }
            }

            Divider().padding(.top, 4)

            HStack {
                Spacer()
                Button(action: {
                    NSApplication.shared.terminate(nil)
                }) {
                    Text(cancelLabel)
                        .frame(maxWidth: .infinity)
                }
                .keyboardShortcut(.cancelAction)
                .frame(width: 100)  // fixed width for Cancel

                Button(action: {
                    var output: [String: Any] = [:]
                    for field in fields {
                        switch field.type {
                        case .text, .texteditor, .filepicker:
                            output[field.bindingKey] = values[field.bindingKey] as? String ?? ""
                        case .checkbox:
                            output[field.bindingKey] = values[field.bindingKey] as? Bool ?? false
                        case .dropdown, .segmented:
                            if let options = field.options, !options.isEmpty {
                                output[field.bindingKey] =
                                    values[field.bindingKey] as? String ?? options.first!
                            }
                        }
                    }
                    if let jsonData = try? JSONSerialization.data(
                        withJSONObject: output, options: []),
                        let jsonString = String(data: jsonData, encoding: .utf8)
                    {
                        print(jsonString)
                    } else {
                        print("Failed to encode JSON.")
                    }
                    NSApplication.shared.terminate(nil)
                }) {
                    Text(okLabel)
                        .frame(maxWidth: .infinity)
                }
                .keyboardShortcut(.defaultAction)
                .frame(width: 100)  // fixed width for OK
            }
        }
        .padding(.vertical, 12)
        .padding(.horizontal, 16)
    }
}

// MARK: - AppDelegate + Setup
class AppDelegate: NSObject, NSApplicationDelegate {
    var window: NSWindow!

    func applicationDidFinishLaunching(_ notification: Notification) {

        // Suppress system debug messages and stderr output
        setenv("OS_ACTIVITY_MODE", "disable", 1)

        // Redirect stderr to /dev/null to suppress system warnings
        freopen("/dev/null", "a", stderr)

        let mainMenu = NSMenu()

        // Application menu (with Quit)
        let appMenuItem = NSMenuItem()
        mainMenu.addItem(appMenuItem)
        let appMenu = NSMenu(title: "Application")
        appMenuItem.submenu = appMenu
        appMenu.addItem(
            withTitle: "Quit", action: #selector(NSApplication.terminate(_:)), keyEquivalent: "q")

        // Edit menu with standard items
        let editMenuItem = NSMenuItem()
        mainMenu.addItem(editMenuItem)
        let editMenu = NSMenu(title: "Edit")
        editMenuItem.submenu = editMenu

        editMenu.addItem(withTitle: "Undo", action: Selector(("undo:")), keyEquivalent: "z")
        editMenu.addItem(withTitle: "Redo", action: Selector(("redo:")), keyEquivalent: "Z")
        editMenu.addItem(NSMenuItem.separator())
        editMenu.addItem(withTitle: "Cut", action: #selector(NSText.cut(_:)), keyEquivalent: "x")
        editMenu.addItem(withTitle: "Copy", action: #selector(NSText.copy(_:)), keyEquivalent: "c")
        editMenu.addItem(
            withTitle: "Paste", action: #selector(NSText.paste(_:)), keyEquivalent: "v")
        editMenu.addItem(
            withTitle: "Delete", action: #selector(NSText.delete(_:)), keyEquivalent: "\u{8}")
        editMenu.addItem(
            withTitle: "Select All", action: #selector(NSText.selectAll(_:)), keyEquivalent: "a")

        NSApplication.shared.mainMenu = mainMenu

        // Use the first command-line argument as config path, if provided.
        let configPath = CommandLine.arguments.count > 1 ? CommandLine.arguments[1] : "config.json"
        let config = loadConfig(from: configPath)

        // Calculate minimum window size based on fields
        let sortedFields = config.sortedFields

        // Calculate height: each field ~40px, texteditor ~120px, spacing + padding + buttons
        var calculatedHeight: CGFloat = 80 // top/bottom padding + button area
        for field in sortedFields {
            switch field.type {
            case .texteditor:
                calculatedHeight += 130 // texteditor height + label + spacing
            default:
                calculatedHeight += 40 // normal field height
            }
            // Add space for note if present
            if let note = field.note, !note.isEmpty {
                calculatedHeight += 20
            }
        }

        // Calculate width: longest label + control width + padding
        let labels = sortedFields.map { $0.label + ":" }
        let font = NSFont.systemFont(ofSize: NSFont.systemFontSize)
        let maxLabelWidth = labels.map { (label: String) -> CGFloat in
            let size = (label as NSString).size(withAttributes: [.font: font])
            return size.width
        }.max() ?? 80
        let calculatedWidth = maxLabelWidth + 400 + 60 // label + control + padding

        let width = config.windowWidth ?? max(500, calculatedWidth)
        let height = config.windowHeight ?? max(300, calculatedHeight)

        let contentView = DynamicDialogView(
            fields: sortedFields,
            okLabel: config.okLabel ?? "OK",
            cancelLabel: config.cancelLabel ?? "Cancel"
        )

        window = NSWindow(
            contentRect: NSRect(x: 0, y: 0, width: width, height: height),
            styleMask: [.titled, .resizable],
            backing: .buffered,
            defer: false
        )
        window.title = config.windowTitle ?? "Dialog"
        window.center()
        window.contentView = NSHostingView(rootView: contentView)

        // Set window level to floating if alwaysOnTop is true
        if config.alwaysOnTop == true {
            window.level = .floating
        }

        window.makeKeyAndOrderFront(nil)
        NSApplication.shared.activate(ignoringOtherApps: true)

        // Set focus to the first text field, if any
        DispatchQueue.main.asyncAfter(deadline: .now() + 0.1) { [weak self] in
            guard let self = self,
                let contentView = self.window.contentView
            else { return }
            if let textField = self.findFirstTextField(in: contentView) {
                self.window.makeFirstResponder(textField)
            }
        }
    }

    func findFirstTextField(in view: NSView?) -> NSView? {
        guard let view = view else { return nil }
        if view is NSTextField { return view }
        for subview in view.subviews {
            if let found = findFirstTextField(in: subview) {
                return found
            }
        }
        return nil
    }
}

// MARK: - Main
let app = NSApplication.shared
let delegate = AppDelegate()
app.delegate = delegate
app.setActivationPolicy(.regular)
app.run()
