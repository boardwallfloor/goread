# Goreader - An Ebook Reader Made with Go

[![GitHub stars](https://img.shields.io/github/stars/goreader/goreader.svg)](https://github.com/boardwallfloor/goread)
[![GitHub forks](https://img.shields.io/github/forks/goreader/goreader.svg)](https://github.com/boardwallfloor/goread/network)
[![GitHub issues](https://img.shields.io/github/issues/goreader/goreader.svg)](https://github.com/boardwallfloor/goread/issues)
[![GitHub license](https://img.shields.io/github/license/goreader/goreader.svg)](https://github.com/boardwallfloor/goread/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/goreader/goreader)](https://goreportcard.com/report/github.com/boardwallfloor/goread)
[![Language](https://img.shields.io/badge/language-Go-blue.svg)](https://golang.org)

Goreader is an ebook reader that has been designed to provide a seamless and enjoyable reading experience for the user. To achieve this, Goreader follows several steps:

1. Specifying the ebook to be read: The user specifies the ebook they want to read, either by selecting it from a library or by uploading it to Goreader.

2. Opening the ebook as a stream: The ebook is opened as a stream, similar to how a zip file is opened. This allows Goreader to access the contents of the ebook without having to fully download it.

3. Determining the mimetype: The mimetype of the ebook is determined to ensure it is an epub file, which is the most common format for ebooks.

4. Locating the content.opf file: The content.opf file is located within the ebook, which contains the spine tag. The spine tag forms the structure of the book, indicating the order in which the various components of the ebook should be displayed.

5. Rendering the book with a given timer: The book is rendered with a given timer to create a readable page as fast as possible. This helps to ensure that the user can start reading the ebook as quickly as possible.

6. Silently rendering more of the book in the background: While the user is reading the ebook, additional pages are silently rendered in the background to improve the reading experience. This helps to reduce the amount of waiting time for the user and makes for a smoother reading experience overall.

By following these steps, Goreader can provide a seamless and enjoyable reading experience for the user.
